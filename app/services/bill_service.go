package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/app/workflows"
	tc "github.com/asheet-bhaskar/billing-service/app/workflows/temporal"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"go.temporal.io/sdk/client"
)

type billService struct {
	repository         repository.BillRepository
	currencyRepository repository.CurrencyRepository
	customerRepository repository.CustomerRepository
	temporalClient     tc.TemporalClient
}

type BillService interface {
	Create(context.Context, *models.BillRequest) (*models.Bill, error)
	GetByID(context.Context, string) (*models.Bill, error)
	AddLineItems(context.Context, *models.LineItem) (*models.LineItem, error)
	RemoveLineItems(context.Context, string, string) (*models.LineItem, error)
	Close(context.Context, string) (*models.Bill, error)
	Invoice(ctx context.Context, billID string) (*models.Invoice, error)
}

func NewBillService(repository repository.BillRepository, currencyRepository repository.CurrencyRepository,
	customerRepository repository.CustomerRepository, temporalClient tc.TemporalClient) BillService {
	return &billService{
		repository:         repository,
		currencyRepository: currencyRepository,
		customerRepository: customerRepository,
		temporalClient:     temporalClient,
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

	workflowID := fmt.Sprintf("BILL-%s", bill.ID)
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "CREATE_BILL_QUEUE",
	}

	_, err = bs.temporalClient.ExecuteWorkflow(context.Background(), options, workflows.BillingWorkflow, bill)
	if err != nil {
		log.Printf("failed to create workflow execution for bill id %s", bill.ID)
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

	lineItem.ID = utils.GetNewUUID()
	lineItem, err = bs.repository.AddLineItems(ctx, lineItem)

	if err != nil {
		log.Printf("error while adding line item %v. error is %s\n", lineItem, err.Error())
		return lineItem, err
	}

	signal := workflows.LineItemSignal{
		BillID: bill.ID,
		ItemID: lineItem.ID,
	}

	err = bs.temporalClient.SignalWorkflow(context.Background(), fmt.Sprintf("BILL-%s", bill.ID), "", "ADD_BILL_ITEM_CHANNEL", signal)
	if err != nil {
		log.Println("Error while signalling the workflow", err)
	}

	return lineItem, nil
}

func (bs *billService) RemoveLineItems(ctx context.Context, billID string, itemID string) (*models.LineItem, error) {
	lineItem, err := bs.repository.GetLineItemByID(ctx, itemID)

	if err == ce.LineItemNotFoundError || err != nil {
		log.Printf("line item not found for id %s\n", itemID)
		return &models.LineItem{}, err
	}

	if lineItem.Removed {
		log.Printf("line item already for id %s\n", itemID)
		return &models.LineItem{}, ce.LineItemAlreadyRemovedError
	}

	bill, err := bs.repository.GetByID(ctx, billID)

	if err == ce.BillNotFoundError || err != nil {
		log.Printf("bill not found for id %s\n", billID)
		return &models.LineItem{}, err
	}

	if bill.Status == "closed" {
		log.Printf("bill is already closed for id %s\n", lineItem.BillID)
		return lineItem, ce.BillClosedError
	}

	lineItemUpdated, err := bs.repository.RemoveLineItems(ctx, lineItem)

	if err != nil {
		log.Printf("error while removing line item %v. error is %s\n", lineItem, err.Error())
		return lineItemUpdated, err
	}

	signal := workflows.LineItemSignal{
		BillID: bill.ID,
		ItemID: lineItemUpdated.ID,
	}

	err = bs.temporalClient.SignalWorkflow(context.Background(), fmt.Sprintf("BILL-%s", bill.ID), "", "REMOVE_BILL_ITEM_CHANNEL", signal)
	if err != nil {
		log.Println("Error while signalling the workflow", err)
	}

	return lineItemUpdated, nil
}

func (bs *billService) Close(ctx context.Context, billID string) (*models.Bill, error) {
	bill, err := bs.repository.GetByID(ctx, billID)

	if err == ce.BillNotFoundError || err != nil {
		log.Printf("bill not found for id %s\n", billID)
		return bill, err
	}

	if bill.Status == "closed" {
		log.Printf("bill is already closed for id %s\n", billID)
		return bill, ce.BillClosedError
	}

	bill, err = bs.repository.Close(ctx, billID)

	if err != nil {
		log.Printf("error while closing bill id %s. error is %s\n", billID, err.Error())
		return bill, err
	}

	return bill, nil
}

func (bs *billService) Invoice(ctx context.Context, billID string) (*models.Invoice, error) {
	invoice := &models.Invoice{}

	bill, err := bs.repository.GetByID(ctx, billID)

	if err == ce.BillNotFoundError || err != nil {
		log.Printf("bill not found for id %s\n", billID)
		return invoice, err
	}

	currency, err := bs.currencyRepository.GetByID(ctx, bill.CurrencyID)
	if err != nil {
		log.Printf("error while fetching currency code for bill id %s\n", billID)
		return invoice, err
	}

	lineItems, err := bs.repository.GetLineItemsByBillID(ctx, bill.ID)
	if err != nil || err != nil {
		log.Printf("error while fetching line items for bill id %s\n", billID)
		return invoice, err
	}

	invoice = models.CreateInvoice(bill, lineItems, currency.Code)

	return invoice, nil
}
