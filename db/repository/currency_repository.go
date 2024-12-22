package repository

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"gorm.io/gorm"
)

type currencyRepository struct {
	db *gorm.DB
}

type CurrencyRepository interface {
	Create(context.Context, *models.Currency) (*models.Currency, error)
	GetByID(context.Context, int64) (*models.Currency, error)
	GetByCode(context.Context, string) (*models.Currency, error)
}

func NewCurrencyRepository(dbClient *gorm.DB) CurrencyRepository {
	return &currencyRepository{
		db: dbClient,
	}
}

func (cr *currencyRepository) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	result := cr.db.Create(&currency)

	if result.Error != nil {
		log.Printf("error occured while creating currency, %v. error is %s", currency, result.Error.Error())
		return currency, result.Error
	}

	return currency, nil
}

func (cr *currencyRepository) GetByID(ctx context.Context, id int64) (*models.Currency, error) {
	currency := &models.Currency{}
	result := cr.db.First(&currency, id)

	if result.Error != nil {
		log.Printf("error occured while querying currency, %d. error is %s", id, result.Error.Error())
		return currency, result.Error
	}

	return currency, nil
}

func (cr *currencyRepository) GetByCode(ctx context.Context, code string) (*models.Currency, error) {
	currency := &models.Currency{}
	result := cr.db.Where("code = ?", code).First(&currency)

	if result.Error != nil {
		log.Printf("error occured while querying currency, %d. error is %s", code, result.Error.Error())
		return currency, result.Error
	}

	return currency, nil
}
