package main

/**
A proxy is an intermediary that acts on behalf of a client or server, relaying requests
and responses between them. It serves as a middle layer between the client
(such as a web browser) and the server (such as a web application).

Forward Proxy:

Used by clients to access servers.
The client sends requests to the proxy, which forwards them to the target server.
Common use cases:
Hiding the client’s identity (IP address).
Bypassing restrictions or filters (e.g., accessing blocked websites).
Caching responses to reduce server load.
Example: A corporate network proxy that filters web traffic.

Reverse Proxy:

Used by servers to manage client requests.
Clients send requests to the reverse proxy, which forwards them to the
appropriate backend server.
Common use cases:
Load balancing: Distributing traffic across multiple servers.
Caching: Storing static resources to reduce load on the server.
Security: Hiding server details from clients and preventing direct access.
SSL/TLS termination: Managing encryption for secure connections.

Key Features of a Proxy
Anonymity: Proxies can hide the client's or server's IP address.
Performance: Proxies can cache data to reduce the need for repeated
requests to the same server.
Security: Proxies can filter malicious requests or prevent unauthorized access.
Load Distribution: Reverse proxies can distribute client requests to multiple
servers, improving scalability.
*/

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Define the backend server to forward requests to
	/**
	Replace http://example.com with the backend server's URL.
	url.Parse parses the backend URL into a format usable by Go's HTTP client.
	*/
	target := "http://example.com"
	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Customize the proxy behavior if needed
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Response status: %s", resp.Status)
		return nil
	}

	/**
	The http.HandleFunc function routes all incoming requests to the reverse proxy.
	proxy.ServeHTTP forwards the request to the backend server
	*/

	// Handle incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	// Start the server
	port := ":8080"
	log.Printf("Reverse proxy server is running on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
