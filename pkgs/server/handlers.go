package server

import (
	"net/http"
	"net/http/httputil"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	targetUrl, err := ExtractBackendUrl(host)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(w, r)
}
