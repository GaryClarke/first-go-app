// File: internal/data/stores.go
package data

import "database/sql"

type Stores struct {
	Books BookStore
}

// NewStores is a constructor function. It takes a database connection
// and returns a Stores struct containing all of our application’s
// data stores (for now, just the BookStore). Using a constructor
// like this keeps the setup logic in one place and makes it easier
// to add more stores later.
func NewStores(db *sql.DB) Stores {
	return Stores{
		Books: BookStore{DB: db},
	}
}
