package database

import "database/sql"

// Permission ...
type Permission struct {
	ID   sql.NullInt64
	Name sql.NullString
}
