package main

import (
	"fmt"
	"net/http"
)

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server 2: %s %s\n", r.Method, r.URL.Path)
}

func main2() {
	http.HandleFunc("/", handler2)
	fmt.Println("Backend Server 2 running on port 8082...")
	http.ListenAndServe(":8082", nil)
}
