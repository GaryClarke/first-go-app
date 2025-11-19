// File: cmd/api/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/garyclarke/first-go-app/internal/data"
	"log"
	"net/http"
)

const version = "1.0.0"

type App struct {
	DB *sql.DB
}

// healthResponse is a struct that represents our JSON response.
// The struct tags (e.g. `json:"status"`) tell the encoder to use lowercase keys in the JSON output.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// The entry point of the Go application.
// This is where the program starts running.
func main() {
	// 1. Open a database connection.
	db, err := data.OpenSQLite()
	if err != nil {
		log.Fatal(err)
	}
	// 2. Close it cleanly when the app shuts down.
	defer db.Close()

	// 3. Migrate and seed
	if err := data.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := data.SeedIfEmpty(db); err != nil {
		log.Fatal(err)
	}

	app := &App{DB: db}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthcheckHandler)
	mux.HandleFunc("GET /books", app.listBooksHandler)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// healthcheckHandler handles incoming requests to /healthz.
// It takes a http.ResponseWriter to write the response,
// and a *http.Request which contains all the request data.
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
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

// writeJSON sends a JSON response to the client.
// It takes a ResponseWriter, a status code, and any value to encode as JSON.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(b)

	return err
}
