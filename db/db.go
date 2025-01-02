package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	gPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// encore:service
type DBClient struct {
	DB *gorm.DB
}

var Clients *DBClient

func InitDBClient(host, port, user, password, name, migrationsPath string) (*DBClient, error) {

	log.Printf("initialising db for %s", name)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	log.Println(dsn)
	dbConn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("failed to create database connection, error %s\n", err.Error())
	}

	gormDB, err := gorm.Open(gPostgres.New(gPostgres.Config{
		Conn: dbConn,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to initialise gorm db client, error %s\n", err.Error())
	}

	Clients = &DBClient{
		DB: gormDB,
	}

	if err = runDBSchemaMigrations(dbConn, migrationsPath); err != nil {
		log.Fatalf("error occured while running database schema migrations, error %s\n", err.Error())
	}

	return Clients, nil
}

func runDBSchemaMigrations(dbConn *sql.DB, migrationsPath string) error {
	log.Println("running database schema migrations")
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsPath), "postgres", driver)

	if err != nil {
		return err
	}
	m.Up()

	log.Println("successfully ran database schema migrations")
	return nil
}
