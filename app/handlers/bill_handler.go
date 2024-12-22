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
func (bs *BillingService) GetBillHandler(ctx context.Context, id int64) (*models.Bill, error) {
	if id <= 0 {
		log.Println("invalid bill id")
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill id",
		}
	}
	bill, err := bs.Bill.GetByID(ctx, id)

	if err == ce.BillNotFoundError {
		log.Printf("bill not found for id %d\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.NotFound,
			Message: "bill not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while fetching bill for id %d\n", id)
		return &models.Bill{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to find bill",
		}
	}

	return bill, nil
}

// encore:api  method=POST path=/bills
func (bs *BillingService) CreateBillHandler(ctx context.Context, request *models.BillRequest) (*models.Bill, error) {
	if !request.IsValid() {
		log.Println("invalid bill request")
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid bill request",
		}
	}

	bill, err := bs.Bill.Create(ctx, request)

	if err == ce.CustomerNotFoundError {
		log.Printf("customer not found for id, %d", request.CustomerID)
		return &models.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: fmt.Sprintf("customer not found for id, %d", request.CustomerID),
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
