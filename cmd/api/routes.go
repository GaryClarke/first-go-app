// File: cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/garyclarke/first-go-app/internal/data"
)

// routes defines the HTTP routes and returns an http.Handler.
//
// In Go, an http.Handler is any type that has a ServeHTTP() method.
// It’s the interface used by the HTTP server to process requests.
//
// Our ServeMux (multiplexer) is a built-in implementation of http.Handler
// that routes requests based on method + path (like "GET /books").
//
// By returning it here, we let main() pass it to http.ListenAndServe,
// which takes over from there and starts handling traffic.
func (app *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", app.healthcheckHandler)
	mux.HandleFunc("GET /books", app.listBooksHandler)
	return mux
}

func (app *App) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Status:  "ok",
		Version: version,
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *App) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := data.GetAll(app.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Write the books to the json response
	if err := writeJSON(w, http.StatusOK, books); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
