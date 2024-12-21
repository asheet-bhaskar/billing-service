package db

import "encore.dev/storage/sqldb"

var db = sqldb.NewDatabase("billing_service", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
