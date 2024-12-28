// //with authorization

// package main

// /**
// encoding/json: For encoding/decoding JSON, useful for handling responses.
// fmt: For formatted I/O (e.g., printing to the console).
// io/ioutil: To read and write from request and response bodies.
// log: For logging errors and messages.
// net/http: The core library for building HTTP servers and clients in Go.
// */
// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// // Route maps endpoint paths to target microservices
// var routes = map[string]string{
// 	"/service-a": "http://localhost:8081",
// 	"/service-b": "http://localhost:8082",
// }

// // ProxyHandler handles incoming requests and forwards them to appropriate microservices
// /**
// Retrieve the Target URL:

// Checks if the requested path exists in the routes map.
// If not, it responds with a 404 - Service not found.
// Forward the Request:

// Uses http.Get(targetURL) to forward the request to the target microservice.
// If an error occurs (e.g., microservice is down), it responds with 500 - Internal Server Error.
// Return the Response:

// Reads the microservice's response body (ioutil.ReadAll).
// Writes the same HTTP status code and response body back to the client.
// */
// func ProxyHandler(w http.ResponseWriter, r *http.Request) {
// 	// Match the request path with the corresponding service
// 	targetURL, exists := routes[r.URL.Path]
// 	if !exists {
// 		http.Error(w, "Service not found", http.StatusNotFound)
// 		return
// 	}

// 	// Forward the request to the target service
// 	resp, err := http.Get(targetURL)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error forwarding request: %s", err.Error()), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Return the response from the microservice
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	w.WriteHeader(resp.StatusCode)
// 	w.Write(body)
// }

// /*
// *
// AuthMiddleware: Ensures that only requests with a valid Authorization header
// are processed.
// http.HandlerFunc: Converts the ProxyHandler function to an http.Handler so it
// can be passed to the middleware.
// http.Handle("/", handler): Uses the middleware-wrapped handler for all incoming
// requests.
// */
// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Authorization")
// 		if token != "valid-token" {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// /** USE this to get authorization tokens
// Invoke-WebRequest -Uri http://localhost:8080/service-a -Headers @{ "Authorization" = "valid-token" } -Method GET
// */

// func main() {

// 	// Wrap the ProxyHandler with the AuthMiddleware
// 	handler := AuthMiddleware(http.HandlerFunc(ProxyHandler))

// 	// Set up HTTP routes
// 	http.Handle("/", handler)

// 	// Start the API Gateway
// 	fmt.Println("API Gateway running on port 8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

//basic

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Route maps endpoint paths to target microservices
var routes = map[string]string{
	"/service-a": "http://localhost:8081",
	"/service-b": "http://localhost:8082",
}

// ProxyHandler handles incoming requests and forwards them to appropriate microservices
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Match the request path with the corresponding service
	targetURL, exists := routes[r.URL.Path]
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Forward the request to the target service
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error forwarding request: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Return the response from the microservice
	body, _ := ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	// Set up HTTP routes
	http.HandleFunc("/", ProxyHandler)

	// Start the API Gateway
	fmt.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// // with rate limiter

// // install go get golang.org/x/time/rate
// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"golang.org/x/time/rate"
// )

// // Route maps endpoint paths to target microservices
// var routes = map[string]string{
// 	"/service-a": "http://localhost:8081",
// 	"/service-b": "http://localhost:8082",
// }

// // ProxyHandler handles incoming requests and forwards them to appropriate microservices
// func ProxyHandler(w http.ResponseWriter, r *http.Request) {
// 	// Match the request path with the corresponding service
// 	targetURL, exists := routes[r.URL.Path]
// 	if !exists {
// 		http.Error(w, "Service not found", http.StatusNotFound)
// 		return
// 	}

// 	// Forward the request to the target service
// 	resp, err := http.Get(targetURL)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error forwarding request: %s", err.Error()), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Return the response from the microservice
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	w.WriteHeader(resp.StatusCode)
// 	w.Write(body)
// }

// // RateLimiter creates a rate limiter middleware
// var limiter = rate.NewLimiter(1, 5) // 1 request per second with a burst of 5

// /*
// *
// The RateLimiterMiddleware checks if the incoming request is allowed by the
// rate limiter (limiter.Allow()).
// If the rate limit has been exceeded, it responds with a 429 Too Many Requests
// status.
// If the request is allowed, it proceeds to the next handler
// (next.ServeHTTP(w, r)).
// */
// func RateLimiterMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !limiter.Allow() {
// 			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// func main() {
// 	// Wrap the ProxyHandler with the RateLimiter middleware
// 	handler := RateLimiterMiddleware(http.HandlerFunc(ProxyHandler))

// 	// Set up HTTP routes
// 	http.Handle("/", handler)

// 	// Start the API Gateway
// 	fmt.Println("API Gateway running on port 8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
