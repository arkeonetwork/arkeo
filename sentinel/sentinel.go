package sentinel

import (
	"arkeo/common"
	"arkeo/sentinel/conf"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

type Proxy struct {
	Metadata   Metadata
	Config     conf.Configuration
	MemStore   *MemStore
	ClaimStore *ClaimStore
}

func NewProxy(config conf.Configuration) Proxy {
	claimStore, err := NewClaimStore(config.ClaimStoreLocation)
	if err != nil {
		panic(err)
	}
	return Proxy{
		Metadata:   NewMetadata(config),
		Config:     config,
		MemStore:   NewMemStore(config.SourceChain),
		ClaimStore: claimStore,
	}
}

// Serve a reverse proxy for a given url
func (p Proxy) serveReverseProxy(w http.ResponseWriter, r *http.Request, host string) {
	// parse the url
	url, _ := url.Parse(fmt.Sprintf("http://%s", host))

	// create the reverse proxy
	proxy := common.NewSingleHostReverseProxy(url)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, r)
}

// Given a request send it to the appropriate url
func (p Proxy) handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	// remove arkauth query arg
	values := r.URL.Query()
	values.Del(QueryArkAuth)
	r.URL.RawQuery = values.Encode()

	parts := strings.Split(r.URL.Path, "/")
	host := parts[1]
	parts = append(parts[:1], parts[1+1:]...)
	r.URL.Path = strings.Join(parts, "/")

	switch host { // nolint
	case "btc-mainnet-fullnode":
		// add username/password to request
		host = fmt.Sprintf("thorchain:password@%s:8332", host)
	}

	p.serveReverseProxy(w, r, host)
}

func (p Proxy) handleMetadata(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	d, _ := json.Marshal(p.Metadata)
	_, _ = w.Write(d)
}

func (p Proxy) handleOpenClaims(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	open_claims := make([]Claim, 0)
	for _, claim := range p.ClaimStore.List() {
		fmt.Printf("Claim: %+v\n", claim)
		if claim.Claimed {
			fmt.Println("already claimed")
			continue
		}
		contract, err := p.MemStore.Get(claim.Key())
		if err != nil {
			fmt.Println("bad fetch")
			continue
		}

		if contract.IsClose(p.MemStore.GetHeight()) {
			_ = p.ClaimStore.Remove(claim.Key()) // clear expired
			fmt.Println("expired")
			continue
		}

		open_claims = append(open_claims, claim)

	}

	d, _ := json.Marshal(open_claims)
	_, _ = w.Write(d)
}

func (p Proxy) handleContract(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "not enough parameters", http.StatusBadRequest)
		return
	}

	providerPK, err := common.NewPubKey(parts[2])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	chain, err := common.NewChain(parts[3])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	spenderPK, err := common.NewPubKey(parts[4])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid spender pubkey", http.StatusBadRequest)
		return
	}

	key := p.MemStore.Key(providerPK.String(), chain.String(), spenderPK.String())
	contract, err := p.MemStore.Get(key)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("fetch contract error: %s", err), http.StatusBadRequest)
		return
	}

	d, _ := json.Marshal(contract)
	_, _ = w.Write(d)
}

func (p Proxy) handleClaim(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "not enough parameters", http.StatusBadRequest)
		return
	}

	providerPK, err := common.NewPubKey(parts[2])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	chain, err := common.NewChain(parts[3])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	spenderPK, err := common.NewPubKey(parts[4])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid spender pubkey", http.StatusBadRequest)
		return
	}

	claim := NewClaim(providerPK, chain, spenderPK, 0, 0, "")
	claim, err = p.ClaimStore.Get(claim.Key())
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("fetch contract error: %s", err), http.StatusBadRequest)
		return
	}

	d, _ := json.Marshal(claim)
	_, _ = w.Write(d)
}

func (p Proxy) Run() {
	log.Println("Starting Sentinel (reverse proxy)....")
	p.Config.Print()

	go p.EventListener(p.Config.EventStreamHost)

	mux := http.NewServeMux()

	// start server
	mux.Handle("/metadata.json", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleMetadata)))
	mux.Handle("/contract/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleContract)))
	mux.Handle("/claim/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleClaim)))
	mux.Handle("/open_claims/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleOpenClaims)))
	mux.Handle("/", p.auth(
		handlers.LoggingHandler(
			os.Stdout,
			handlers.ProxyHeaders(
				http.HandlerFunc(p.handleRequestAndRedirect),
			),
		),
	))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", p.Config.Port), mux))
}
