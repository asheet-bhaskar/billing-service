package db

import (
	"log"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type DBClient struct {
	DB     *gorm.DB
	TestDB *gorm.DB
}

var db = sqldb.NewDatabase("billing_service", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

var testDB = sqldb.NewDatabase("billing_service_test", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

var Clients *DBClient

func InitDBClient() (*DBClient, error) {
	log.Println("initialising database clients for application databse")
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Stdlib(),
	}))
	if err != nil {
		log.Fatalf("failed to initialise application database, error", err.Error())
	}

	log.Println("initialising database clients for application tests")
	testDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: testDB.Stdlib(),
	}))
	if err != nil {
		log.Fatalf("failed to initialise test database, error", err.Error())
	}

	Clients = &DBClient{
		DB:     db,
		TestDB: testDB,
	}
	return Clients, nil
}
