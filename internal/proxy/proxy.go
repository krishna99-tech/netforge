package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Start starts a reverse proxy on the given port, forwarding requests to the target URL.
func Start(port int, target string) error {
	targetURL, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Optional: Log requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[PROXY] %s %s -> %s\n", r.Method, r.URL.Path, target)
		proxy.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting reverse proxy on %s, forwarding to %s...\n", addr, target)
	return http.ListenAndServe(addr, nil)
}
