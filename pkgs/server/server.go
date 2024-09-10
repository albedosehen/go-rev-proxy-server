package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/acme/autocert"
)

var DOMAINS = strings.Split(os.Getenv("DOMAINS"), ",")
var BACKENDS = make(map[string]string)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(BACKENDS) != 0 {
		return
	}

	for _, domain := range DOMAINS {
		subdomain := strings.Split(domain, ".")[0]
		envVarKey := strings.ToUpper(subdomain) + "_BACKEND"

		backendURL := os.Getenv(envVarKey)
		if backendURL == "" {
			log.Fatalf("Backend URL for domain %s is not defined. Expected environment variable: %s", domain, envVarKey)
		}

		BACKENDS[domain] = backendURL
	}
}

func ExtractBackendUrl(host string) (*url.URL, error) {
	if backendURL, exists := BACKENDS[host]; exists {
		// Parse the backend URL
		targetURL, err := url.Parse(backendURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing backend URL: %v", err)
		}

		return targetURL, nil
	}

	return nil, fmt.Errorf("backend URL for host %s not found", host)
}

func StartServer(useHTTPRedirect bool) {
	Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", ProxyHandler)

	// Let's Encrypt autocert manager
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(DOMAINS...),             // Use the domains from the .env file
		Cache:      autocert.DirCache(os.Getenv("CERT_CACHE_DIR")), // Use the certificate cache directory from the .env file
	}

	// HTTP server for Let's Encrypt challenge on port 80
	if useHTTPRedirect {
		go func() {
			log.Fatal(http.ListenAndServe(":80", certManager.HTTPHandler(nil)))
		}()
	} else {
		go func() {
			log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
			})))
		}()
	}

	// HTTPS server using Let's Encrypt certificates
	server := &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: certManager.TLSConfig(),
	}

	// Start the HTTPS server
	log.Println("Starting server on :443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
