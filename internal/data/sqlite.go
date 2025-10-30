// File: internal/data/sqlite.go
package data

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite" // Blank import: registers the "sqlite" driver with database/sql
	"time"
)

// DSN (Data Source Name) tells SQLite where/how to store the database.
//
// Here we’re using a file called books.db in the project root.
// The ?_pragma=busy_timeout(5000) part tells SQLite to wait up to 5 seconds
// if the database is locked, instead of failing immediately. This helps avoid
// “database is locked” errors when we do quick consecutive writes in demos.
const dsn = "file:books.db?_pragma=busy_timeout(5000)"

// OpenSQLite opens a database connection pool for SQLite and checks it works.
//
// A *sql.DB is not a single connection. It’s a pool of connections managed
// by the database/sql package. With SQLite we restrict this pool to 1
// connection (since SQLite only allows one writer at a time).
func OpenSQLite() (*sql.DB, error) {
	// sql.Open doesn’t actually establish any connections yet.
	// It just prepares the pool with the driver and DSN.
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Limit the pool so SQLite doesn’t trip over concurrency.
	db.SetMaxOpenConns(1)    // at most 1 open connection at a time
	db.SetMaxIdleConns(1)    // keep at most 1 idle connection ready
	db.SetConnMaxLifetime(0) // don’t recycle connections by age

	// To *actually* verify the connection, we ping it.
	// PingContext tries to open a connection from the pool and check it works.
	//
	// We wrap it in a context with a timeout. A context carries deadlines,
	// cancellation signals, and request-scoped values across API boundaries.
	// Here it ensures Ping won’t hang forever — it will fail if it takes longer
	// than 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		// If ping fails, close the pool before returning the error.
		_ = db.Close()
		return nil, err
	}

	// At this point we know the DSN is valid and SQLite is ready to use.
	return db, nil
}
