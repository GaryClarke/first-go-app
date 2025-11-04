// File: internal/data/migrate.go
package data

import "database/sql"

// Migrate creates the tables we need for this project.
// In real-world projects, migrations are usually run with a separate tool,
// but for this course we run them at startup to keep things simple.
func Migrate(db *sql.DB) error {
	// DDL = Data Definition Language. This is the subset of SQL used to define
	// or modify the database schema (tables, columns, indexes, constraints).
	// It's different from DML (Data Manipulation Language), which you use to
	// insert, update, delete, or query rows of data.
	const ddl = `
CREATE TABLE IF NOT EXISTS books (
  id     INTEGER PRIMARY KEY AUTOINCREMENT,
  title  TEXT NOT NULL,
  author TEXT,
  year   INTEGER
);`
	// Exec runs the DDL statement. If the table already exists, the
	// IF NOT EXISTS clause ensures nothing bad happens.
	_, err := db.Exec(ddl)
	return err
}

// SeedIfEmpty inserts demo books if the table is empty.
// This is just for demo purposes - we’ll remove it later once POST /books is implemented.
// Once added..run the following commands:
// go run ./cmd/api
// sqlite3 books.db "SELECT * FROM books;"
// You should see two rows in the output. If you run it again, no duplicates will be added.
func SeedIfEmpty(db *sql.DB) error {
	var count int

	// Count how many rows currently exist in the books table.
	err := db.QueryRow(`SELECT COUNT(*) FROM books`).
		Scan(&count)
	if err != nil {
		return err
	}

	// If there are already rows, do nothing.
	if count > 0 {
		return nil
	}

	// Insert two demo books.
	_, err = db.Exec(`
INSERT INTO books (title, author, year) VALUES
  ('The Go Programming Language', 'Alan Donovan', 2015),
  ('Designing Data-Intensive Applications', 'Martin Kleppmann', 2017)`)

	return err
}
