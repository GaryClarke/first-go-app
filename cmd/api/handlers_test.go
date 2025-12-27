// File: cmd/api/handlers_test.go
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/garyclarke/first-go-app/internal/data"
	"net/http"
	"net/http/httptest"

	// Blank import: registers the "sqlite" driver with database/sql
	// The blank identifier (_) tells Go we're importing this package only for its side effect
	// (the driver registration), not to use any of its functions or types directly
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestApp(t *testing.T) *App {
	// Mark this function as a test helper
	// This tells Go's test runner that if a test fails, the error should point
	// to the actual test function that called setupTestApp, not to this helper function
	t.Helper()

	// Open an in-memory SQLite database for testing
	// ":memory:" is a special SQLite connection string that creates a database
	// that only exists in RAM - it's perfect for tests because it's fast and
	// automatically cleaned up when the connection closes
	db, err := sql.Open("sqlite", ":memory:")

	// If we can't open the database, the test can't continue
	// t.Fatal() stops the test immediately and reports the error
	if err != nil {
		t.Fatal(err)
	}

	// Register a cleanup function that will run after the test finishes
	// This ensures the database connection is properly closed, even if the test fails
	// t.Cleanup() runs in reverse order (last registered, first executed) after the test ends
	// This is better than defer because defer runs when the function returns,
	// but cleanup runs after the entire test function completes
	t.Cleanup(func() {
		db.Close()
	})

	// Run the database migrations to create the tables we need
	// If migration fails, we can't run the test, so we stop immediately
	if err := data.Migrate(db); err != nil {
		t.Fatal(err)
	}

	// Seed the database with initial test data if it's empty
	// This gives us a known starting state for our tests
	// If seeding fails, we stop the test
	if err := data.SeedIfEmpty(db); err != nil {
		t.Fatal(err)
	}

	// Return a new App instance with the test database
	// This is what our test handlers will use instead of the real database
	return &App{Stores: data.NewStores(db)}
}

func TestListBooksHandler(t *testing.T) {
	// setup test
	app := setupTestApp(t)

	// create test request
	req := httptest.NewRequest(http.MethodGet, "/books", http.NoBody)

	// create test recorder
	rr := httptest.NewRecorder()

	// invoke the handler
	app.listBooksHandler(rr, req)

	// check status code
	if rr.Code != http.StatusOK {
		t.Errorf("want status code %d; got %d", http.StatusOK, rr.Code)
	}

	// create a bookResponse var
	var resp bookResponse

	// decode the response body into the booksResponse var
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	// check length of books
	booksCount := len(resp.Books)
	if booksCount != 2 {
		t.Errorf("want books count of 2; got %d", booksCount)
	}
}

func TestShowBookHandler(t *testing.T) {
	// setup test
	app := setupTestApp(t)

	// create test request
	req := httptest.NewRequest(http.MethodGet, "/books/1", http.NoBody)

	// create test recorder
	rr := httptest.NewRecorder()

	// invoke the handler
	app.showBookHandler(rr, req)

	// check status code
	if rr.Code != http.StatusOK { // 200
		t.Errorf("want status code %d; got %d", http.StatusOK, rr.Code)
	}

	// create a book var
	var book data.Book

	// decode the response body into the book var
	if err := json.NewDecoder(rr.Body).Decode(&book); err != nil {
		t.Fatal(err)
	}

	// expected book

	// check book against expected
}
