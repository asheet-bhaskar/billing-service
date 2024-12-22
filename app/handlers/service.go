package handlers

import (
	service "github.com/asheet-bhaskar/billing-service/app/services"
	"github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/db/repository"
)

// encore:service
type BillingService struct {
	Bill     service.BillService
	Customer service.CustomerService
	Currency service.CurrencyService
}

func initBillingService() (*BillingService, error) {
	dbClient, _ := db.InitDBClient()
	BillRepo := repository.NewBillRepository(dbClient.DB)
	CustomerRepo := repository.NewCustomerRepository(dbClient.DB)
	CurrencyRepo := repository.NewCurrencyRepository(dbClient.DB)

	return &BillingService{
		Bill:     service.NewBillService(BillRepo, CurrencyRepo, CustomerRepo),
		Customer: service.NewCustomerService(CustomerRepo),
		Currency: service.NewCurrencyService(CurrencyRepo),
	}, nil
}
