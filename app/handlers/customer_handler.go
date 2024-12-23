package handlers

import (
	"context"
	"log"

	"encore.dev/beta/errs"
	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
)

// encore:api method=GET path=/customers/:id
func (bs *BillingService) GetCustomerHandler(ctx context.Context, id string) (*models.Customer, error) {
	if id == "" {
		log.Println("invalid customer id")
		return &models.Customer{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid customer id",
		}
	}
	customer, err := bs.Customer.GetByID(ctx, id)

	if err == ce.CustomerNotFoundError {
		log.Printf("customer not found for id %s\n", id)
		return &models.Customer{}, &errs.Error{
			Code:    errs.NotFound,
			Message: "customer not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while fetching customer for id %s\n", id)
		return &models.Customer{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to get customer",
		}
	}

	return customer, nil
}

// encore:api  method=POST path=/customers
func (bs *BillingService) CreateCustomerHandler(ctx context.Context, request *models.Customer) (*models.Customer, error) {
	if !request.IsValid() {
		log.Println("invalid customer request")
		return &models.Customer{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid customer request",
		}
	}

	customer, err := bs.Customer.Create(ctx, request)

	if err == ce.CustomerAlreadyExistError {
		log.Println("customer already exists")
		return &models.Customer{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "customer already exists",
		}
	}

	if err != nil {
		log.Println("failed to create customer")
		return &models.Customer{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to create customer",
		}
	}

	return customer, nil
}
