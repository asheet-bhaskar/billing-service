package service

import (
	"context"
	"log"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
)

type billService struct {
	repository         repository.BillRepository
	currencyRepository repository.CurrencyRepository
	customerRepository repository.CustomerRepository
}

type BillService interface {
	Create(context.Context, *models.BillRequest) (*models.Bill, error)
	GetByID(context.Context, string) (*models.Bill, error)
	AddLineItems(context.Context, *models.LineItem) (*models.LineItem, error)
	RemoveLineItems(context.Context, *models.LineItem) (*models.LineItem, error)
	Close(context.Context, string) (*models.Bill, error)
}

func NewBillService(repository repository.BillRepository, currencyRepository repository.CurrencyRepository, customerRepository repository.CustomerRepository) BillService {
	return &billService{
		repository:         repository,
		currencyRepository: currencyRepository,
		customerRepository: customerRepository,
	}
}

func (bs *billService) Create(ctx context.Context, request *models.BillRequest) (*models.Bill, error) {
	currency, err := bs.currencyRepository.GetByCode(ctx, request.CurrencyCode)
	if err != nil {
		log.Printf("error while finding the currency for code %s\n", request.CurrencyCode)
		return &models.Bill{}, err
	}

	customer, err := bs.customerRepository.GetByID(ctx, request.CustomerID)

	if err != nil {
		log.Printf("error while finding the customer for id %s\n", request.CustomerID)
		return &models.Bill{}, err
	}
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: request.Description,
		CustomerID:  customer.ID,
		CurrencyID:  currency.ID,
		Status:      "open",
		TotalAmount: 0.0,
		PeriodStart: request.PeriodStart,
		PeriodEnd:   request.PeriodEnd,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	bill, err = bs.repository.Create(ctx, bill)
	if err != nil {
		log.Printf("error occured while creating bill. error %s\n", err.Error())
		return &models.Bill{}, err
	}

	return bill, nil
}

func (bs *billService) GetByID(ctx context.Context, id string) (*models.Bill, error) {
	bill, err := bs.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("error occured while fetching bill with id %s. error %s\n", id, err.Error())
		return &models.Bill{}, err
	}

	return bill, nil
}

func (bs *billService) AddLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	bill, err := bs.repository.GetByID(ctx, lineItem.BillID)

	if err == ce.BillNotFoundError || err != nil {
		log.Printf("bill not found for id %s\n", lineItem.BillID)
		return lineItem, err
	}

	if bill.Status == "closed" {
		log.Printf("bill is already closed for id %s\n", lineItem.BillID)
		return lineItem, ce.BillClosedError
	}

	lineItem, err = bs.repository.AddLineItems(ctx, lineItem)

	if err != nil {
		log.Printf("error while adding line item %v. error is %s\n", lineItem, err.Error())
		return lineItem, err
	}

	return lineItem, nil
}

func (bs *billService) RemoveLineItems(context.Context, *models.LineItem) (*models.LineItem, error) {
	return &models.LineItem{}, nil
}

func (bs *billService) Close(context.Context, string) (*models.Bill, error) {
	return &models.Bill{}, nil
}

func (bs *billService) Invoice(context.Context, string) (*models.Invoice, error) {
	return &models.Invoice{}, nil
}
