package main

/**
fmt: Used to format and print messages, such as logging the server being selected.
log: Used for logging information, such as errors and server selections.
http: Provides functions to create and handle HTTP requests and responses.
httputil: Contains utilities like the NewSingleHostReverseProxy function,
which helps forward requests to backend servers.
url: Provides utilities for URL parsing, which is used when defining the backend servers.
sync: Provides the Mutex type, which is used to safely manage concurrent access to shared
resources (like the round-robin index).
*/
import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	servers []string
	mu      sync.Mutex
	index   int
}

/*
*
NewLoadBalancer: A constructor function that initializes and returns a new LoadBalancer
object with the provided list of servers.
*/
func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{servers: servers}
}

// GetNextServer returns the next backend server in a round-robin fashion
/**
GetNextServer: This function returns the next backend server in the list, using
round-robin logic:
lb.mu.Lock(): Locks the Mutex to ensure only one goroutine can access the index at a time.

defer lb.mu.Unlock(): Ensures that the lock is released after the function completes.

lb.index: Selects the current server using the index.

(lb.index + 1) % len(lb.servers): Increments the index and wraps it around to the start

of the list when reaching the end (round-robin behavior).
Logging: The selected server is logged for debugging purposes.
*/
func (lb *LoadBalancer) GetNextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers) // Round-robin logic

	// Log the server being used for debugging
	log.Printf("Selecting backend server: %s\n", server)

	return server
}

// ProxyHandler proxies the incoming request to the backend server
/**
ProxyHandler: This function handles HTTP requests coming to the load balancer and forwards
them to the appropriate backend server.

GetNextServer(): Calls the function we defined earlier to get the next server in the
round-robin rotation.

url.Parse(server): Parses the backend server URL so that we can create a reverse proxy.

httputil.NewSingleHostReverseProxy(url): Creates a reverse proxy that forwards the request
to the selected backend server.

proxy.ServeHTTP(w, r): This function actually proxies the incoming request (r) to the
backend server and returns the response to the client.
*/
func (lb *LoadBalancer) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the next server
	server := lb.GetNextServer()

	// Log the server selection (for debugging)
	log.Printf("Forwarding request to: %s\n", server)

	// Create a URL for the proxy server
	url, err := url.Parse(server)
	if err != nil {
		http.Error(w, "Failed to parse backend server URL", http.StatusInternalServerError)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Proxy the request to the backend server
	proxy.ServeHTTP(w, r)
}

func main() {
	// List of backend servers
	backendServers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	// Create a new load balancer
	lb := NewLoadBalancer(backendServers)

	// Start the load balancer server
	http.HandleFunc("/", lb.ProxyHandler)

	// Run the load balancer on port 8080
	fmt.Println("Load Balancer running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
