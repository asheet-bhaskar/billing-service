package service

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
)

type customerService struct {
	repository repository.CustomerRepository
}

type CustomerService interface {
	Create(context.Context, *models.Customer) (*models.Customer, error)
	GetByID(context.Context, int64) (*models.Customer, error)
}

func NewCustomerService(repository repository.CustomerRepository) CustomerService {
	return &customerService{
		repository: repository,
	}
}

func (cs *customerService) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer, err := cs.repository.Create(ctx, customer)
	if err != nil {
		log.Printf("error occured while creating custome. error %s\n", err.Error())
		return &models.Customer{}, err
	}

	return customer, nil
}

func (cs *customerService) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	customer, err := cs.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("error occured while fetching customer with id %d. error %s\n", id, err.Error())
		return &models.Customer{}, err
	}

	return customer, nil
}
