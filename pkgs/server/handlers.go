package server

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: Method=%s, Host=%s, URL=%s, RemoteAddr=%s", r.Method, r.Host, r.URL.String(), r.RemoteAddr)

	host := r.Host
	targetUrl, err := ExtractBackendUrl(host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(w, r)
}
