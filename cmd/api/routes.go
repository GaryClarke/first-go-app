// File: cmd/api/routes.go
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/garyclarke/first-go-app/internal/request"
	"net/http"
	"strconv"

	"github.com/garyclarke/first-go-app/internal/data"
)

type bookResponse struct {
	Books []data.Book `json:"books"`
}

// healthResponse is a struct that represents our JSON response.
// The struct tags (e.g. `json:"status"`) tell the encoder to use lowercase keys in the JSON output.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

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
	mux.HandleFunc("GET /books/{id}", app.showBookHandler)
	mux.HandleFunc("POST /books", app.createBookHandler)
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
	books, err := app.Stores.Books.GetAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := bookResponse{Books: books}

	// Write the books to the json response
	if err := writeJSON(w, http.StatusOK, resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *App) showBookHandler(w http.ResponseWriter, r *http.Request) {
	// Get the value of id
	idString := r.PathValue("id")
	// Convert to an int for the db lookup
	id, err := strconv.ParseInt(idString, 10, 64)
	// Validate the id
	if err != nil || id < 1 {
		// Return not found if can't be validated
		http.NotFound(w, r)
		return
	}

	book, err := app.Stores.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			http.NotFound(w, r) // 404
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Write the json response
	if err := writeJSON(w, http.StatusOK, book); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *App) createBookHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Declare an input struct to hold the incoming JSON data.
	var input request.FullBookRequest

	// Step 2: Decode the request body into the input struct.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Step 3: Validate the input data.
	// Coming soon: title required, year in range, etc.

	// Step 4: Create a Book struct with the validated data.
	// For now, we'll hard-code a fake book to simulate this.
	book := &data.Book{
		ID:     3, // fake ID
		Title:  "The Go Workshop",
		Author: "Delio D'Anna",
		Year:   2021,
	}

	// Step 5: Return the created book as JSON with a 201 Created status.
	if err := writeJSON(w, http.StatusCreated, book); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
