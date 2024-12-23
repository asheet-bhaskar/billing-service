package handlers

import (
	"context"
	"fmt"
	"log"

	"encore.dev/beta/errs"
	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
)

// encore:api method=GET path=/bills/:id
func (bs *APIService) GetBillHandler(ctx context.Context, id string) (*models.Bill, error) {
	if id == "" {
		log.Println("invalid bill id")
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill id",
		}
	}
	bill, err := bs.Bill.GetByID(ctx, id)

	if err == ce.BillNotFoundError {
		log.Printf("bill not found for id %s\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.NotFound,
			Message: "bill not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while fetching bill for is %d\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to find bill",
		}
	}

	return bill, nil
}

// encore:api  method=POST path=/bills
func (bs *APIService) CreateBillHandler(ctx context.Context, request *models.BillRequest) (*models.Bill, error) {
	if !request.IsValid() {
		log.Println("invalid bill request")
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill request",
		}
	}

	bill, err := bs.Bill.Create(ctx, request)

	if err == ce.CustomerNotFoundError {
		log.Printf("customer not found for id, %s", request.CustomerID)
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: fmt.Sprintf("customer not found for id, %s", request.CustomerID),
		}
	}

	if err == ce.CurrencyNotFoundError {
		log.Printf("currency not found for code, %s", request.CurrencyCode)
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: fmt.Sprintf("currency not found for code, %s", request.CurrencyCode),
		}
	}

	if err == ce.BillAlreadyExistError {
		log.Println("bill already exists")
		return &models.Bill{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "bill already exists",
		}
	}

	if err != nil {
		log.Println("failed to create bill")
		return &models.Bill{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to create bill",
		}
	}

	return bill, nil
}

//encore:api method=POST path=/bills/items
func (bs *APIService) AddLineItemsHandler(ctx context.Context, lineItem models.LineItem) (*models.LineItem, error) {
	if lineItem.BillID == "" {
		log.Println("invalid bill id")
		return &lineItem, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill id",
		}
	}

	item, err := bs.Bill.AddLineItems(ctx, &lineItem)

	if err == ce.BillNotFoundError {
		log.Println("bill not found")
		return &lineItem, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill not found",
		}
	}

	if err == ce.BillClosedError {
		log.Println("bill closed already")
		return &lineItem, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill closed already",
		}
	}

	if err != nil {
		log.Println("failed to add line item")
		return &lineItem, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to add line item",
		}
	}

	return item, nil
}

//encore:api method=PUT path=/bills/:billID/items/:itemID
func (bs *APIService) RemoveLineItemsHandler(ctx context.Context, billID string, itemID string) (*models.LineItem, error) {
	if billID == "" || itemID == "" {
		log.Println("invalid bill id or item id")
		return &models.LineItem{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill idor item id",
		}
	}

	item, err := bs.Bill.RemoveLineItems(ctx, billID, itemID)

	if err == ce.LineItemAlreadyRemovedError {
		log.Println("line item aleady removed")
		return item, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "line item aleady removed",
		}
	}

	if err == ce.BillNotFoundError {
		log.Println("bill not found")
		return item, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill not found",
		}
	}

	if err == ce.BillClosedError {
		log.Println("bill closed already")
		return item, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill closed already",
		}
	}

	if err != nil {
		log.Println("failed to remove line item")
		return item, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to remove line item",
		}
	}

	return item, nil
}

// encore:api method=GET path=/bills/:id/invoice
func (bs *APIService) GetInvoiceHandler(ctx context.Context, id string) (*models.Invoice, error) {
	if id == "" {
		log.Println("invalid bill id")
		return &models.Invoice{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill id",
		}
	}
	invoice, err := bs.Bill.Invoice(ctx, id)

	if err == ce.BillNotFoundError {
		log.Printf("bill not found for id %s\n", id)
		return &models.Invoice{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill not found",
		}
	}

	if err == ce.CurrencyNotFoundError {
		log.Printf("currency not found for id %s\n", id)
		return &models.Invoice{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "currency not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while fetching bill for is %s\n", id)
		return &models.Invoice{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to find bill",
		}
	}

	return invoice, nil
}

// encore:api method=PUT path=/bills/:id/close
func (bs *APIService) CloseBillHandler(ctx context.Context, id string) (*models.Bill, error) {
	if id == "" {
		log.Println("invalid bill id")
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill id",
		}
	}
	bill, err := bs.Bill.Close(ctx, id)

	if err == ce.BillNotFoundError {
		log.Printf("bill not found for id %s\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.NotFound,
			Message: "bill not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while closing bill for is %d\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to find bill",
		}
	}

	return bill, nil
}
