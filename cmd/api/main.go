// File: cmd/api/main.go
package main

import (
	"encoding/json"
	"github.com/garyclarke/first-go-app/internal/data"
	"log"
	"net/http"
)

const version = "1.0.0"

// App holds the dependencies for our HTTP handlers.
// Instead of passing a raw *sql.DB around, we now store
// a data.Stores value. This gives our handlers access to
// all the application’s data stores (currently just Books)
// through a single field.
type App struct {
	Stores data.Stores
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

	// Build our App with all its dependencies.
	// For now this means the data stores, created from the DB connection.
	app := &App{Stores: data.NewStores(db)}

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", app.routes()); err != nil {
		log.Fatal(err)
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
