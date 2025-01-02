package handlers

import (
	"log"

	"encore.dev/config"
	appConfig "github.com/asheet-bhaskar/billing-service/app/config"
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

var applicationConfig = config.Load[appConfig.Config]()

func initAPIService() (*APIService, error) {
	dbClient, _ := db.InitDBClient(applicationConfig)
	BillRepo := repository.NewBillRepository(dbClient.DB)
	CustomerRepo := repository.NewCustomerRepository(dbClient.DB)
	CurrencyRepo := repository.NewCurrencyRepository(dbClient.DB)
	temporalClient, err := client.NewClient(client.Options{
		HostPort:  applicationConfig.TemporalHostPort(),
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
