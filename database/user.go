package database

import "database/sql"

// User ...
type User struct {
	ID       sql.NullInt64
	Username sql.NullString
	Password sql.NullString
	FullName sql.NullString
	Age      sql.NullInt64
	Address  sql.NullString
	GroupID  sql.NullInt64
}

// UserPermissions ...
type UserPermissions struct {
	ID           sql.NullInt64
	UserID       sql.NullInt64
	PermissionID sql.NullInt64
}

// UserAPI ...
type UserAPI struct {
	Username  string
	Password  string
	FullName  string
	Age       int64
	Address   string
	GroupName string
}
