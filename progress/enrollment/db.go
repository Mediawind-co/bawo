package enrollment

import "encore.dev/storage/sqldb"

// Database for the enrollment service.
var db = sqldb.NewDatabase("enrollment", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
