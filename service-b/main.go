package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Response from Service B")
	})
	fmt.Println("Service B running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
