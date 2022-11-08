package switchd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"mercury/switch/conf"

	"github.com/gorilla/handlers"
)

type Proxy struct {
	Metadata Metadata
	Config   conf.Configuration
}

func NewProxy() Proxy {
	config := conf.NewConfiguration()
	return Proxy{
		Metadata: NewMetadata(config),
		Config:   config,
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
	d, _ := json.Marshal(p.Metadata)
	w.Write(d)
}

func (p Proxy) Run() {
	log.Printf("Starting server.... :%s", p.Config.Port)

	mux := http.NewServeMux()

	// start server
	mux.Handle("/", handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(http.HandlerFunc(p.handleRequestAndRedirect))))
	mux.Handle("/metadata.json", handlers.LoggingHandler(os.Stdout, enforceJSONHandler(http.HandlerFunc(p.handleMetadata))))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", p.Config.Port), mux))
}
