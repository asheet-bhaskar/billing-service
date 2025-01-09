package repository

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"gorm.io/gorm"
)

type currencyRepository struct {
	db *gorm.DB
}

type CurrencyRepository interface {
	Create(context.Context, *models.Currency) (*models.Currency, error)
	GetByID(context.Context, string) (*models.Currency, error)
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

func (cr *currencyRepository) GetByID(ctx context.Context, id string) (*models.Currency, error) {
	currency := &models.Currency{}
	result := cr.db.Where("id = ?", id).First(&currency)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("currency not found for id %s\n", id)
		return currency, ce.CurrencyNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while querying currency, %s. error is %s", id, result.Error.Error())
		return currency, result.Error
	}

	return currency, nil
}

func (cr *currencyRepository) GetByCode(ctx context.Context, code string) (*models.Currency, error) {
	currency := &models.Currency{}
	result := cr.db.Where("code = ?", code).Find(&currency)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("currency not found for code %s\n", code)
		return currency, ce.CurrencyNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while querying currency, %s. error is %s", code, result.Error.Error())
		return currency, result.Error
	}

	return currency, nil
}
