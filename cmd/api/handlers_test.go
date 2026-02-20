// File: cmd/api/handlers_test.go
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/garyclarke/first-go-app/internal/data"

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
	app.routes().ServeHTTP(rr, req)

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

	// send the request through the router
	// We’re updating our tests to send requests through the router, just like real HTTP traffic would.
	// This ensures path parameters like {id} are parsed correctly, and keeps all of our handler tests consistent.
	app.routes().ServeHTTP(rr, req)

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
	expected := data.Book{
		ID:     1,
		Title:  "The Go Programming Language",
		Author: "Alan Donovan",
		Year:   2015,
	}

	// check book against expected
	if book != expected {
		t.Errorf("want %#v; got %#v", expected, book)
	}
}

func TestCreateBookHandler_ValidInput(t *testing.T) {
	// Setup
	app := setupTestApp(t)

	// Create a JSON payload as a raw string.
	// The handler reads the request body via an io.Reader (e.g. when decoding JSON).
	// strings.NewReader wraps our string, so it can be read the same way: it implements io.Reader.
	// So we're not sending a real HTTP request; we're just giving the handler something that
	// behaves like a body (readable bytes) so we can test it without a real server or network.
	body := strings.NewReader(`{
		"title":"Testing Go",
		"author":"Gary Clarke",
		"year":2030
	}`)

	// Make POST request with valid JSON
	req := httptest.NewRequest(http.MethodPost, "/books", body)

	// It's important to set the Content-Type header so the server knows to treat the body as JSON.
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder — this captures the response that would be sent to a client.
	rr := httptest.NewRecorder()

	// Route the request through our app’s router, just like the real server would.
	app.routes().ServeHTTP(rr, req)

	// Assert: 201 status, correct JSON response
	if rr.Code != http.StatusCreated {
		t.Errorf("want status code %d; got %d", http.StatusCreated, rr.Code)
	}

	// decode the response body into the book var
	var book data.Book

	if err := json.NewDecoder(rr.Body).Decode(&book); err != nil {
		t.Fatal(err)
	}

	if book.ID < 1 {
		t.Errorf("expected book to have a positive value ID")
	}
	if book.Title != "Testing Go" {
		t.Errorf("expected title to be 'Testing Go'; got %q", book.Title)
	}
	if book.Author != "Gary Clarke" {
		t.Errorf("expected author to be 'Gary Clarke'; got %q", book.Author)
	}
	if book.Year != 2030 {
		t.Errorf("expected year to be 2030; got %d", book.Year)
	}

	// Verify book exists in the DB
	stored, err := app.Stores.Books.Get(book.ID)
	if err != nil {
		t.Fatalf("failed to fetch book from DB: %v", err)
	}

	// verify stored book matches returned book
	// stored is a *Book (a pointer), but book is a value.
	// To compare them properly, we dereference stored using *stored
	// so we’re comparing two Book values directly.
	if *stored != book {
		t.Errorf("book in DB does not match response. got: %#v", stored)
	}
}

func TestCreateBookHandler_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		wantCode int
		wantKeys []string // expected keys in the "errors" object of the response
	}{
		{
			name:     "missing all fields",
			payload:  `{}`,
			wantCode: http.StatusUnprocessableEntity,
			wantKeys: []string{"title", "author", "year"},
		},
		{
			name:     "missing title",
			payload:  `{"author": "Gary", "year": 2023}`,
			wantCode: http.StatusUnprocessableEntity,
			wantKeys: []string{"title"},
		},
		{
			name:     "missing author",
			payload:  `{"title": "Testing Go", "year": 2023}`,
			wantCode: http.StatusUnprocessableEntity,
			wantKeys: []string{"author"},
		},
		{
			name:     "invalid year (zero)",
			payload:  `{"title": "Testing Go", "author": "Gary", "year": 0}`,
			wantCode: http.StatusUnprocessableEntity,
			wantKeys: []string{"year"},
		},
		{
			name:     "invalid JSON format",
			payload:  `{`,
			wantCode: http.StatusBadRequest,
			wantKeys: nil, // No "errors" object expected — it's a decoding error
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Setup a fresh in-memory app
			app := setupTestApp(t)

			// Step 2: Create a request using tc.payload as the body
			body := strings.NewReader(tc.payload)

			req := httptest.NewRequest(http.MethodPost, "/books", body)

			// Step 3: Set the Content-Type header to application/json
			req.Header.Set("Content-Type", "application/json")

			// Step 4: Create a response recorder
			rr := httptest.NewRecorder()

			// Step 5: Send the request through the app router
			app.routes().ServeHTTP(rr, req)

			// Step 6: Assert that the status code matches tc.wantCode
			if rr.Code != tc.wantCode {
				t.Errorf("want status code %d; got %d", tc.wantCode, rr.Code)
			}

			// Step 7: If tc.wantKeys is not nil,
			//         decode the response JSON into a map[string]any
			//         and check that all expected error keys exist
			if tc.wantKeys != nil {
				var resp map[string]any
				err := json.NewDecoder(rr.Body).Decode(&resp)
				if err != nil {
					t.Fatal(err)
				}
				// Assert that "errors" field exists and is a map[string]any
				errorsMap, ok := resp["errors"].(map[string]any)
				if !ok {
					t.Fatalf("expected 'errors' field in response, got: %#v", resp)
				}

				// Check that all expected error keys exist
				for _, key := range tc.wantKeys {
					if _, ok := errorsMap[key]; !ok {
						t.Errorf("expected error for key %q in response", key)
					}
				}
			}
		})
	}
}
