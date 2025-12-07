// File: internal/data/books.go
package data

import (
	"context"
	"database/sql"
	"time"
)

// BookStore wraps a sql.DB connection pool.
// It provides methods for working with books in the database
// (for example, fetching all books or looking up a book by ID).
type BookStore struct {
	DB *sql.DB
}

func GetAll(db *sql.DB) ([]Book, error) {
	// Define the SQL query to fetch all books, ordered by ID
	const query = `SELECT id, title, author, year FROM books ORDER BY id`

	// Create a context with a 3-second timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Ensure the context is cleaned up when this function exits (defer)
	defer cancel()

	// Execute the query using the context (will timeout after 3 seconds if not done)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// Close the database rows when we're done reading them
	defer rows.Close()

	// Create an empty slice to hold our Book structs
	var books []Book

	// Loop through each row returned from the database
	for rows.Next() {
		// Create a new Book struct for this row
		var b Book
		// Scan the row's columns into the Book struct fields
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year); err != nil {
			return nil, err
		}
		// Add this book to our books slice
		books = append(books, b)
	}

	// Check if there were any errors during iteration (not caught by rows.Scan)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of books and nil for no error
	return books, nil
}
