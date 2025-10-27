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
	// Create the response data using our struct and constant
	response := healthResponse{
		Status:  "ok",
		Version: version,
	}

	// Attempt to write the JSON response using our writeJSON helper.
	// If something goes wrong (e.g. the data can't be encoded),
	// we return a 500 Internal Server Error to the client.
	if err := writeJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// writeJSON sends a JSON response to the client.
// It takes a ResponseWriter, a status code, and any value to encode as JSON.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	// Step 1: Convert the value (v) into a JSON byte slice
	b, err := json.Marshal(v)
	if err != nil {
		// If the JSON encoding fails, return the error
		return err
	}

	// Step 2: Set the Content-Type header so the client knows we're sending JSON
	w.Header().Set("Content-Type", "application/json")

	// Step 3: Set the HTTP status code (e.g. 200, 400, 500)
	w.WriteHeader(status)

	// Step 4: Write the JSON bytes to the response body
	_, err = w.Write(b)

	// Return any error from writing to the response
	return err
}
