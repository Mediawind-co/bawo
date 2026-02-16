package content

import "encore.dev/storage/sqldb"

// Database for the content service.
var db = sqldb.NewDatabase("content", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
