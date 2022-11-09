package sentinel

import (
	"encoding/json"
	"fmt"
	"log"
	"mercury/common"
	"mercury/sentinel/conf"
	"net/http"
	"net/http/httputil"
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
func (p Proxy) serveReverseProxy(w http.ResponseWriter, r *http.Request) {
	// parse the url
	url, _ := url.Parse(p.Config.ProxyHost)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, r)
}

// Given a request send it to the appropriate url
func (p Proxy) handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	p.serveReverseProxy(w, r)
}

func (p Proxy) handleMetadata(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	d, _ := json.Marshal(p.Metadata)
	_, _ = w.Write(d)
}

func (p Proxy) handleContract(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "not enough parameters", http.StatusBadRequest)
		return
	}

	providerPK, err := common.NewPubKey(parts[1])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	chain, err := common.NewChain(parts[2])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	spenderPK, err := common.NewPubKey(parts[3])
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
	if len(parts) < 4 {
		http.Error(w, "not enough parameters", http.StatusBadRequest)
		return
	}

	providerPK, err := common.NewPubKey(parts[1])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	chain, err := common.NewChain(parts[2])
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("bad provider pubkey: %s", err), http.StatusBadRequest)
		return
	}

	spenderPK, err := common.NewPubKey(parts[3])
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

	mux := http.NewServeMux()

	// start server
	mux.Handle("/metadata.json", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleMetadata)))
	mux.Handle("/contract/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleContract)))
	mux.Handle("/claim/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handleClaim)))
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
