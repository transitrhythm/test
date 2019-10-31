package main

import (
	"net/http"
    "github.com/r3labs/sse"
)

func main() {
	server := sse.New()

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/events", server.HTTPHandler)

	http.ListenAndServe(":8080", mux)
}