package tracker

import "encore.dev/storage/sqldb"

// Database for the tracker service.
var db = sqldb.NewDatabase("tracker", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
