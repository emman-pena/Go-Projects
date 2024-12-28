// proxy.go
package main

/**
log: This package is used to log messages and errors.
net/http: Provides HTTP client and server functionality.
net/http/httputil: Contains utilities for HTTP, such as
ReverseProxy which helps in proxying requests to other servers.
net/url: Provides URL parsing and manipulation functions.
strings: Provides functions for string manipulation
(used here for checking URL paths).
*/
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// url.Parse to create URL objects from strings.
func main() {
	// Define the upstream services
	service1URL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal("Error parsing Service 1 URL:", err)
	}
	service2URL, err := url.Parse("http://localhost:8082")
	if err != nil {
		log.Fatal("Error parsing Service 2 URL:", err)
	}

	// Create reverse proxies for both services

	service1Proxy := httputil.NewSingleHostReverseProxy(service1URL)
	service2Proxy := httputil.NewSingleHostReverseProxy(service2URL)

	// Handle routing based on URL path
	/**
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)):
	This line tells the server to handle all incoming HTTP requests with
	the function provided.

	w is the http.ResponseWriter, used to write responses.
	r is the http.Request, which contains the incoming request data,
	such as the URL path.
	strings.HasPrefix(r.URL.Path, "/service1"): We check if the URL path
	starts with /service1. This is how we route the traffic to Service 1.

	If the path starts with /service1, we call service1Proxy.ServeHTTP(w, r),
	which forwards the request to service1.
	strings.HasPrefix(r.URL.Path, "/service2"): Similarly, if the path starts
	with /service2, the request is forwarded to service2.

	http.NotFound(w, r): If the path doesn't match /service1 or /service2,
	we return a 404 error indicating that the requested resource was not found.
	*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/service1") {
			service1Proxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/service2") {
			service2Proxy.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	/**
	http.ListenAndServe(":8080", nil): This starts an HTTP server that listens
	on port 8080. The second argument is nil, meaning we’re using the default
	http.ServeMux multiplexer (the router).
	*/
	log.Println("Proxy is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
