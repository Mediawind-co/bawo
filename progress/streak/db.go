package streak

import "encore.dev/storage/sqldb"

// Database for the streak service.
var db = sqldb.NewDatabase("streak", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
