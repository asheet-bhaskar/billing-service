package repository

import (
	"context"
	"log"

	models "github.com/asheet-bhaskar/billing-service/app/services/customer"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

type customerRepository interface {
	Create(context.Context, *models.Customer) (*models.Customer, error)
	GetByID(context.Context, int64) (*models.Customer, error)
}

func NewCustomerRepository(dbClient *gorm.DB) customerRepository {
	return &CustomerRepository{
		db: dbClient,
	}
}

func (cr *CustomerRepository) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	result := cr.db.Create(&customer)

	if result.Error != nil {
		log.Printf("error occured while creating customer, %v. error is %s", customer, result.Error.Error())
		return customer, result.Error
	}

	return customer, nil
}

func (cr *CustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	customer := &models.Customer{}
	result := cr.db.First(&customer, id)

	if result.Error != nil {
		log.Printf("error occured while querying customer, %d. error is %s", id, result.Error.Error())
		return customer, result.Error
	}

	return customer, nil
}
