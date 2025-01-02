package db

import (
	"fmt"
	"log"

	"encore.dev/storage/sqldb"
	"github.com/asheet-bhaskar/billing-service/app/config"
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

func InitDBClient(appConfig config.Config) (*DBClient, error) {
	log.Println("initialising database clients for application databse")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		appConfig.DBHost(), appConfig.DBUser(), appConfig.DBPassword(), appConfig.DBName(), appConfig.DBPort())

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
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
