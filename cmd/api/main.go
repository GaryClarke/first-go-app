// File: cmd/api/main.go
package main

import (
	"encoding/json"
	"net/http"
)

const version = "1.0.0"

// healthResponse is a struct that represents our JSON response.
// The struct tags (e.g. `json:"status"`) tell the encoder to use lowercase keys in the JSON output.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

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
	// Tell the client that we're sending JSON back in the response
	w.Header().Set("Content-Type", "application/json")

	// Create the response data using our struct and constant
	response := healthResponse{
		Status:  "ok",
		Version: version,
	}

	// Convert the map to JSON and write it to the response
	json.NewEncoder(w).Encode(response)
}
