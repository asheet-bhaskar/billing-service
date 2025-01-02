package handlers

import (
	"log"

	"encore.dev/config"
	service "github.com/asheet-bhaskar/billing-service/app/services"
	"github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	"github.com/asheet-bhaskar/billing-service/worker"
	"go.temporal.io/sdk/client"
)

// encore:service
type APIService struct {
	Bill     service.BillService
	Customer service.CustomerService
	Currency service.CurrencyService
}

type Config struct {
	TemporalHostPort       config.String
	DBHost                 config.String
	DBPort                 config.String
	DBUser                 config.String
	DBPassword             config.String
	DBName                 config.String
	DBSchemaMigrationsPath config.String
}

var appConfig = config.Load[Config]()

func initAPIService() (*APIService, error) {
	dbClient, _ := db.InitDBClient(appConfig.DBHost(), appConfig.DBPort(), appConfig.DBUser(), appConfig.DBUser(), appConfig.DBName(), appConfig.DBSchemaMigrationsPath())
	BillRepo := repository.NewBillRepository(dbClient.DB)
	CustomerRepo := repository.NewCustomerRepository(dbClient.DB)
	CurrencyRepo := repository.NewCurrencyRepository(dbClient.DB)
	temporalClient, err := client.NewClient(client.Options{
		HostPort:  appConfig.TemporalHostPort(),
		Namespace: "default",
	})

	if err != nil {
		log.Fatal("Failed to initiate temporal client")
	}

	log.Println("starting temporal worker")
	go worker.Start(temporalClient)

	return &APIService{
		Bill:     service.NewBillService(BillRepo, CurrencyRepo, CustomerRepo, temporalClient),
		Customer: service.NewCustomerService(CustomerRepo),
		Currency: service.NewCurrencyService(CurrencyRepo),
	}, nil
}
