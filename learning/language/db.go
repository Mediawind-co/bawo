package language

import "encore.dev/storage/sqldb"

// Database for the language service.
var db = sqldb.NewDatabase("language", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
