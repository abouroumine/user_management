package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver     = "mysql"
	user       = "user"
	password   = "pass"
	host       = "localhost"
	port       = "3306"
	dbName     = "userManagement"
	connection = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?parseTime=true"
)

// ConnectMySQL ...
func ConnectMySQL() (*sql.DB, error) {
	return sql.Open(driver, connection)
}
