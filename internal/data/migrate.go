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
`
	// Exec runs the DDL statement. If the table already exists, the
	// IF NOT EXISTS clause ensures nothing bad happens.
	_, err := db.Exec(ddl)
	return err
}

