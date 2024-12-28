package main

import (
	"fmt"
	"net/http"
)

/**
backendHandler1: This function handles requests for Backend 1 and sends a response
that indicates the request was handled by Server 1.
http.ListenAndServe(":8081", nil): Starts the backend server on port 8081.
*/

func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server 1: %s %s\n", r.Method, r.URL.Path)
}

func main1() {
	http.HandleFunc("/", handler1)
	fmt.Println("Backend Server 1 running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
