package database

import "database/sql"

// Group ...
type Group struct {
	ID   sql.NullInt64
	Name sql.NullString
}
