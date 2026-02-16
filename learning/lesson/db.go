package lesson

import "encore.dev/storage/sqldb"

// Database for the lesson service.
var db = sqldb.NewDatabase("lesson", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
