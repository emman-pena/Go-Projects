package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Response from Service A")
	})
	fmt.Println("Service A running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
