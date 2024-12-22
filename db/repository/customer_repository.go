package repository

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

type CustomerRepository interface {
	Create(context.Context, *models.Customer) (*models.Customer, error)
	GetByID(context.Context, int64) (*models.Customer, error)
}

func NewCustomerRepository(dbClient *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: dbClient,
	}
}

func (cr *customerRepository) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	result := cr.db.Create(&customer)

	if result.Error != nil {
		log.Printf("error occured while creating customer, %v. error is %s", customer, result.Error.Error())
		return customer, result.Error
	}

	return customer, nil
}

func (cr *customerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	customer := &models.Customer{}
	result := cr.db.First(&customer, id)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("customer not found for id %d\n", id)
		return customer, ce.CustomerNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while querying customer, %d. error is %s", id, result.Error.Error())
		return customer, result.Error
	}

	return customer, nil
}
