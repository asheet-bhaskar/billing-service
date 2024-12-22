package service

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
)

type currencyService struct {
	repository repository.CurrencyRepository
}

type CurrencyService interface {
	Create(context.Context, *models.Currency) (*models.Currency, error)
	GetByID(context.Context, int64) (*models.Currency, error)
}

func NewCurrencyService(repository repository.CurrencyRepository) CurrencyService {
	return &currencyService{
		repository: repository,
	}
}

func (cs *currencyService) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	currency, err := cs.repository.Create(ctx, currency)
	if err != nil {
		log.Printf("error occured while creating currency. error %s\n", err.Error())
		return &models.Currency{}, err
	}

	return currency, nil
}

func (cs *currencyService) GetByID(ctx context.Context, id int64) (*models.Currency, error) {
	currency, err := cs.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("error occured while fetching currency with id %d. error %s\n", id, err.Error())
		return &models.Currency{}, err
	}

	return currency, nil
}
