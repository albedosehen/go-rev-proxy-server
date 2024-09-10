package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

var (
	DOMAINS  []string
	BACKENDS = make(map[string]string)
)

func Init() {
	if len(DOMAINS) > 0 && len(BACKENDS) > 0 {
		// Domains and backends are already initialized
		return
	}

	domainsEnv := os.Getenv("DOMAINS")
	if domainsEnv == "" {
		log.Fatal("DOMAINS environment variable is not set")
	}

	// Populate DOMAINS slice
	DOMAINS = strings.Split(domainsEnv, ",")

	// Populate BACKENDS map
	for _, domain := range DOMAINS {
		subdomain := strings.Split(domain, ".")[0]
		envVarKey := strings.ToUpper(subdomain) + "_BACKEND" // e.g www.example.com is the domain, then the environment variable should be WWW_BACKEND

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

func StartHttpServerInBackground(useStandardHTTPChallengeHandling bool, certManager *autocert.Manager) {
	if useStandardHTTPChallengeHandling {
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
}

func StartServer(useStandardHTTPChallengeHandling bool) {
	// TODO: This will be a required manual call with optional parameters to initialize the server configuration in the future
	Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", ProxyHandler)

	// Let's Encrypt autocert manager
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(DOMAINS...),             // Use the domains from the .env file
		Cache:      autocert.DirCache(os.Getenv("CERT_CACHE_DIR")), // Use the certificate cache directory from the .env file
	}

	StartHttpServerInBackground(useStandardHTTPChallengeHandling, &certManager)

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
