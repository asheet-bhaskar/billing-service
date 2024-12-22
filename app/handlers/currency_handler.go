package handlers

import (
	"context"
	"log"

	"encore.dev/beta/errs"
	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
)

// encore:api method=GET path=/currencies/:id
func (bs *BillingService) GetCurrencyHandler(ctx context.Context, id int64) (*models.Currency, error) {
	if id <= 0 {
		log.Println("invalid currency id")
		return &models.Currency{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid currency id",
		}
	}
	currency, err := bs.Currency.GetByID(ctx, id)

	if err == ce.CurrencyNotFoundError {
		log.Printf("currency not found for id %d\n", id)
		return &models.Currency{}, &errs.Error{
			Code:    errs.NotFound,
			Message: "currency not found",
		}
	}

	if err != nil {
		log.Printf("error occurred while fetching currency for id %d\n", id)
		return &models.Currency{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to get currency",
		}
	}

	return currency, nil
}

// encore:api  method=POST path=/currencies
func (bs *BillingService) CreateCurrencyHandler(ctx context.Context, request *models.Currency) (*models.Currency, error) {
	if !request.IsValid() {
		log.Println("invalid currency request")
		return &models.Currency{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid currency request",
		}
	}

	currency, err := bs.Currency.Create(ctx, request)

	if err == ce.CurrencyAlreadyExistError {
		log.Println("currency already exists")
		return &models.Currency{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "currency already exists",
		}
	}

	if err != nil {
		log.Println("failed to create currency")
		return &models.Currency{}, &errs.Error{
			Code:    errs.Unknown,
			Message: "failed to create currency",
		}
	}

	return currency, nil
}
