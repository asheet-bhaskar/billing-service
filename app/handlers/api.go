package handlers

import (
	"log"

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

func initAPIService() (*APIService, error) {
	dbClient, _ := db.InitDBClient()
	BillRepo := repository.NewBillRepository(dbClient.DB)
	CustomerRepo := repository.NewCustomerRepository(dbClient.DB)
	CurrencyRepo := repository.NewCurrencyRepository(dbClient.DB)
	temporalClient, err := client.NewClient(client.Options{})

	if err != nil {
		log.Fatal("Failed to start temporal")
	}

	log.Println("starting temporal worker")
	go worker.Start(temporalClient)

	return &APIService{
		Bill:     service.NewBillService(BillRepo, CurrencyRepo, CustomerRepo, temporalClient),
		Customer: service.NewCustomerService(CustomerRepo),
		Currency: service.NewCurrencyService(CurrencyRepo),
	}, nil
}
