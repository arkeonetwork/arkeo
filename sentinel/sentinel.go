package sentinel

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	"github.com/sirupsen/logrus"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
)

type Proxy struct {
	Metadata            Metadata
	Config              conf.Configuration
	MemStore            *MemStore
	ClaimStore          *ClaimStore
	ContractConfigStore *ContractConfigurationStore
	ProviderConfigStore *ProviderConfigurationStore
	logger              log.Logger
	proxies             map[string]*url.URL
}

func NewProxy(config conf.Configuration) (Proxy, error) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	claimStore, err := NewClaimStore(config.ClaimStoreLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create claim store with error: %s", err))
		return Proxy{}, fmt.Errorf("failed to create claim store with error: %s", err)
	}
	contractConfigStore, err := NewContractConfigurationStore(config.ContractConfigStoreLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create contract config store with error: %s", err))
		return Proxy{}, fmt.Errorf("failed to create contract config store with error: %s", err)
	}

	providerConfigStore, err := NewProviderConfigurationStore(config.ProviderConfigStoreLocation)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create provider config store with error: %s", err))
		return Proxy{}, fmt.Errorf("failed to create provider config store with error: %s", err)
	}

	return Proxy{
		Metadata:            NewMetadata(config),
		Config:              config,
		MemStore:            NewMemStore(config.SourceChain, logger),
		ClaimStore:          claimStore,
		ContractConfigStore: contractConfigStore,
		proxies:             loadProxies(),
		logger:              logger,
		ProviderConfigStore: providerConfigStore,
	}, nil
}

func loadProxies() map[string]*url.URL {
	proxies := make(map[string]*url.URL)
	for serviceName := range common.ServiceLookup {
		// if we have an override for a given service, parse that instead of
		// the default below
		env, envOk := os.LookupEnv(strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")))
		if envOk {
			proxies[serviceName] = common.MustParseURL(env)
			continue
		}

		// parse default values for services
		switch serviceName {
		case "btc-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://infra:password@bitcoin-daemon:8332")
		case "bch-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://infra:password@bitcoin-cash-daemon:8332")
		case "doge-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://infra:password@doge-daemon:8332")
		case "ltc-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://infra:password@litecoin-daemon:8332")
		case "arkeo-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://arkeo:1317")
		case "eth-mainnet-fullnode":
			proxies[serviceName] = common.MustParseURL("http://ethereum-daemon:8545")
		case "gaia-mainnet-grpc":
			proxies[serviceName] = common.MustParseURL("http://gaia-daemon:9090")
		case "gaia-mainnet-rpc":
			proxies[serviceName] = common.MustParseURL("http://gaia-daemon:26657")
		case "mock":
			proxies[serviceName] = common.MustParseURL("http://localhost:3765")
		default:
			proxies[serviceName] = common.MustParseURL(fmt.Sprintf("http://%s", serviceName))
		}
	}
	return proxies
}

// Given a request send it to the appropriate url
func (p Proxy) handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	// Limit the Size of incoming requests

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // TODO: Check

	// remove arkauth query arg
	values := r.URL.Query()
	values.Del(QueryArkAuth)
	r.URL.RawQuery = values.Encode()

	parts := strings.Split(r.URL.Path, "/")

	serviceName := r.Header.Get(ServiceHeader)
	pulledFromPath := false
	if len(serviceName) == 0 && len(parts) > 1 {
		pulledFromPath = true
		serviceName = parts[1]
	}

	uri, exists := p.proxies[serviceName]
	if !exists {
		respondWithError(w, "could not find service", http.StatusBadRequest)
		return
	}

	r.URL.Scheme = uri.Scheme
	r.URL.Host = uri.Host
	r.URL.User = uri.User
	if pulledFromPath {
		parts[1] = uri.Path // replace service name with uri path (if exists)
		r.URL.Path = path.Join(parts...)
	}

	// Sanitize URL
	// ensure path always has "/" prefix
	if len(r.URL.Path) > 1 && !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = fmt.Sprintf("/%s", r.URL.Path)
	}

	// check for the WebSocket upgrade header
	if websocket.IsWebSocketUpgrade(r) {
		fmt.Println(">>>>>>> Entering websocket....")
		wsProxyURL := *r.URL
		wsProxyURL.Scheme = "ws" // use the WebSocket scheme
		wsProxy := websocketproxy.NewProxy(&wsProxyURL)
		wsProxy.ServeHTTP(w, r)
		return
	}

	// Serve a reverse proxy for a given url
	// create the reverse proxy
	proxy := common.NewSingleHostReverseProxy(r.URL)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, r)
}

func (p Proxy) handleMetadata(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	d, _ := json.Marshal(p.Metadata)
	_, _ = w.Write(d)
}

func (p Proxy) handleContract(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		respondWithError(w, "missing id in uri", http.StatusBadRequest)
		return
	}

	contractId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		p.logger.Error("fail to parse contract id", "error", err, "id", id)
		respondWithError(w, fmt.Sprintf("bad contract id: %s", err), http.StatusBadRequest)
		return
	}

	conf, err := p.ContractConfigStore.Get(contractId)
	if err != nil {
		p.logger.Error("fail to fetch contract", "error", err, "id", contractId)
		respondWithError(w, fmt.Sprintf("bad contract id: %s", err), http.StatusBadRequest)
		return
	}

	// check authorization
	var auth ContractAuth
	raw := r.Header.Get(QueryContract)
	if len(raw) > 0 {
		auth, err = parseContractAuth(raw)
		if err != nil {
			p.logger.Error("fail to parse contract auth", "error", err, "auth", raw[0])
			respondWithError(w, fmt.Sprintf("bad contract auth: %s", err), http.StatusBadRequest)
			return
		}

		contract, err := p.MemStore.Get(conf.Key())
		if err != nil {
			p.logger.Error("fail to fetch contract", "error", err, "id", conf.Key())
			respondWithError(w, fmt.Sprintf("missing contract: %s", err), http.StatusNotFound)
			return
		}
		if err := auth.Validate(conf.LastTimeStamp, contract.Client); err != nil {
			p.logger.Error("fail to validate contract auth", "error", err, "auth", auth.String())
			respondWithError(w, fmt.Sprintf("bad contract auth: %s", err), http.StatusBadRequest)
			return
		}
		conf.LastTimeStamp = auth.Timestamp
		if err := p.ContractConfigStore.Set(conf); err != nil {
			p.logger.Error("fail to save contract config", "error", err, "auth", auth.String())
			respondWithError(w, fmt.Sprintf("fail to save contract config: %s", err), http.StatusBadRequest)
			return
		}
	} else {
		p.logger.Error("missing contract auth")
		respondWithError(w, fmt.Sprintf("missing contract auth: %s", err), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		d, _ := json.Marshal(conf)
		_, _ = w.Write(d)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		type PostContractConfig struct {
			PerUserRateLimit     int      `json:"per_user_rate_limit"`
			CORs                 CORs     `json:"cors"`
			WhitelistIPAddresses []string `json:"white_listed_ip_addresses"`
		}
		var changes PostContractConfig
		if err := json.Unmarshal(body, &changes); err != nil {
			http.Error(w, "Error unmarshaling JSON data", http.StatusBadRequest)
			return
		}

		conf.PerUserRateLimit = changes.PerUserRateLimit
		conf.CORs = changes.CORs
		conf.WhitelistIPAddresses = changes.WhitelistIPAddresses
		err = p.ContractConfigStore.Set(conf)
		if err != nil {
			p.logger.Error("fail to save contract config", "error", err, "id", conf.ContractId)
			respondWithError(w, fmt.Sprintf("failed to save contract config: %s", err), http.StatusInternalServerError)
			return
		}
	default:
		p.logger.Error("unsupported request method", "method", r.Method)
		respondWithError(w, fmt.Sprintf("unsupported request method: %s", r.Method), http.StatusBadRequest)
	}
}

func (p Proxy) handleOpenClaims(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	open_claims := make([]Claim, 0)
	for _, claim := range p.ClaimStore.List() {
		if claim.Claimed {
			p.logger.Info("already claimed")
			continue
		}
		contract, err := p.MemStore.Get(claim.Key())
		if err != nil {
			p.logger.Error("bad fetch")
			continue
		}

		if contract.IsExpired(p.MemStore.GetHeight()) {
			_ = p.ClaimStore.Remove(claim.Key()) // clear expired
			p.logger.Info("claim expired")
			continue
		}

		open_claims = append(open_claims, claim)
	}

	d, _ := json.Marshal(open_claims)
	_, _ = w.Write(d)
}

func (p Proxy) handleActiveContract(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	service, ok := vars["service"]
	if !ok {
		respondWithError(w, "missing service in uri", http.StatusBadRequest)
		return
	}
	pubkey, ok := vars["spender"]
	if !ok {
		respondWithError(w, "missing spender pubkey in uri", http.StatusBadRequest)
		return
	}

	providerPK := p.Config.ProviderPubKey

	r.URL.Path = fmt.Sprintf("/arkeo/active-contract/%s/%s/%s",
		providerPK.String(),
		service,
		pubkey,
	)
	proxy := httputil.NewSingleHostReverseProxy(p.proxies["arkeo-mainnet-fullnode"])
	proxy.ServeHTTP(w, r)
}

func (p Proxy) handleClaim(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		respondWithError(w, "missing id in uri", http.StatusBadRequest)
		return
	}
	contractId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		p.logger.Error("fail to parse contractId", "error", err, "contractId", id)
		respondWithError(w, fmt.Sprintf("bad contractId: %s", err), http.StatusBadRequest)
		return
	}

	claim := NewClaim(contractId, nil, 0, "")
	claim, err = p.ClaimStore.Get(claim.Key())
	p.logger.Info(fmt.Sprintf("claim data %v", claim))
	if err != nil {
		p.logger.Error("fail to get claim from memstore", "error", err, "key", claim.Key())
		respondWithError(w, fmt.Sprintf("fetch contract error: %s", err), http.StatusBadRequest)
		return
	}

	d, _ := json.Marshal(claim)
	_, _ = w.Write(d)
}

func (p Proxy) Run() {
	p.logger.Info("Starting Sentinel (reverse proxy)....")
	p.Config.Print()

	go p.EventListener(p.Config.EventStreamHost)

	router := p.getRouter()

	// Configure Logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Add the Logrus middleware to the router
	loggingRouter := p.logrusMiddleware(router)

	// Check if TLS certificates are configured
	if p.Config.TLS.HasTLS() {
		// Start a goroutine that listens on port 80 and redirects HTTP to HTTPS
		go func() {
			redirectServer := &http.Server{
				Addr: fmt.Sprintf(":%s", p.Config.Port),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
				}),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
				IdleTimeout:  5 * time.Second,
			}
			if err := redirectServer.ListenAndServe(); err != nil {
				panic(err)
			}
		}()

		// Start HTTPS server on port 443
		server := &http.Server{
			Addr:              ":443",
			Handler:           loggingRouter,
			ReadTimeout:       5 * time.Second, // TODO: updated it to use config
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       120 * time.Second,
			TLSConfig: &tls.Config{
				// Policies
				MinVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_AES_128_GCM_SHA256,
					tls.TLS_AES_256_GCM_SHA384,
					tls.TLS_CHACHA20_POLY1305_SHA256,
				},
			},
			MaxHeaderBytes: 1 << 20,
		}
		if err := server.ListenAndServeTLS(p.Config.TLS.Cert, p.Config.TLS.Key); err != nil {
			panic(err)
		}
	} else {
		// Start HTTP server on the configured port
		server := &http.Server{
			Addr:              fmt.Sprintf(":%s", p.Config.Port),
			Handler:           loggingRouter,
			ReadTimeout:       5 * time.Second, // TODO: updated it to use config
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20,
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}

func (p *Proxy) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(RoutesMetaData, http.HandlerFunc(p.handleMetadata)).Methods(http.MethodGet)
	router.HandleFunc(RoutesActiveContract, http.HandlerFunc(p.handleActiveContract)).Methods(http.MethodGet)
	router.HandleFunc(RoutesClaim, http.HandlerFunc(p.handleClaim)).Methods(http.MethodGet)
	router.HandleFunc(RoutesOpenClaims, http.HandlerFunc(p.handleOpenClaims)).Methods(http.MethodGet)
	router.HandleFunc(RouteManage, http.HandlerFunc(p.handleContract)).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(RouteProviderData, http.HandlerFunc(p.handleProviderData)).Methods(http.MethodGet)
	router.PathPrefix("/").Handler(
		p.auth(
			handlers.ProxyHeaders(
				http.HandlerFunc(p.handleRequestAndRedirect),
			),
		),
	)
	return router
}

func (p *Proxy) logrusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"remote": p.getRemoteAddr(r),
		})

		logger.Info("New request")
		next.ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		response = []byte(`{"error": "failed to marshal response payload"}`)
		code = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func (p Proxy) handleProviderData(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	serviceString, ok := vars["service"]
	if !ok {
		respondWithError(w, "missing service in uri", http.StatusBadRequest)
		return
	}
	service := common.Service(common.ServiceLookup[serviceString])

	providerConfigData, err := p.ProviderConfigStore.Get(p.Config.ProviderPubKey, service.String())
	if err != nil {
		p.logger.Error("failed to get provider details", "error", err, "provider", p.Config.ProviderPubKey)
		respondWithError(w, fmt.Sprintf("Invalid Provider: %s", err), http.StatusBadRequest)
		return
	}

	d, _ := json.Marshal(providerConfigData)
	_, _ = w.Write(d)
}
