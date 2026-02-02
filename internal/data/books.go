// File: internal/data/books.go
package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// BookStore wraps a sql.DB connection pool.
// It provides methods for working with books in the database
// (for example, fetching all books or looking up a book by ID).
type BookStore struct {
	DB *sql.DB
}

func (s *BookStore) GetAll() ([]Book, error) {
	// Define the SQL query to fetch all books, ordered by ID
	query := `SELECT id, title, author, year FROM books ORDER BY id`

	// Create a context with a 3-second timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Ensure the context is cleaned up when this function exits (defer)
	defer cancel()

	// Execute the query using the context (will timeout after 3 seconds if not done)
	rows, err := s.DB.QueryContext(ctx, query)
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

func (s *BookStore) Get(id int64) (*Book, error) {
	// In SQLite, auto-incremented IDs start at 1.
	// To avoid making a pointless database query,
	// we immediately return sql.ErrNoRows if the ID is less than 1.
	if id < 1 {
		// we reuse sql.ErrNoRows, which is exactly what the database driver gives us when a query finds nothing.
		// That way our handler only needs one simple check: was the error sql.ErrNoRows? If yes, return 404
		return nil, sql.ErrNoRows
	}

	query := `SELECT id, title, author, year FROM books WHERE id = ?`

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Declare a Book struct to hold the data returned by the query.
	var book Book

	// Query and scan into book
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Year,
	)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookStore) Insert(book *Book) (*Book, error) {
	// query
	query := `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`
	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute query
	res, err := s.DB.ExecContext(ctx, query, book.Title, book.Author, book.Year)
	if err != nil {
		return nil, err
	}
	// get the id
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// set id on book
	book.ID = id

	log.Printf("Inserted book: %+v", book)
	// return the book
	return book, nil
}
