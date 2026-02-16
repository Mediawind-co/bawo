package user

import (
	"encore.dev/storage/sqldb"
)

// Database for the user service
var db = sqldb.NewDatabase("users", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
