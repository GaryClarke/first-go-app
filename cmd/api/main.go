// File: cmd/api/main.go
package main

import (
	"net/http"
)

// The entry point of the Go application.
// This is where the program starts running.
func main() {
	// Create a new HTTP request multiplexer (router).
	// This will match incoming requests to the correct handler functions.
	mux := http.NewServeMux()

	// Register a handler function for GET requests to the /healthz endpoint.
	// When a GET request hits /healthz, the healthcheckHandler function will be called.
	mux.HandleFunc("GET /healthz", healthcheckHandler)

	// Start the HTTP server on port 8080 and pass in the mux (router) to handle requests.
	// This call blocks - the program runs until the server is stopped.
	http.ListenAndServe(":8080", mux)
}

// healthcheckHandler handles incoming requests to /healthz.
// It takes a http.ResponseWriter to write the response,
// and a *http.Request which contains all the request data.
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}
