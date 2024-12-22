package service

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
)

type billService struct {
	repository repository.BillRepository
}

type BillService interface {
	Create(context.Context, *models.Bill) (*models.Bill, error)
	GetByID(context.Context, int64) (*models.Bill, error)
}

func NewBillService(repository repository.BillRepository) BillService {
	return &billService{
		repository: repository,
	}
}

func (cs *billService) Create(ctx context.Context, bill *models.Bill) (*models.Bill, error) {
	bill, err := cs.repository.Create(ctx, bill)
	if err != nil {
		log.Printf("error occured while creating bill. error %s\n", err.Error())
		return &models.Bill{}, err
	}

	return bill, nil
}

func (cs *billService) GetByID(ctx context.Context, id int64) (*models.Bill, error) {
	bill, err := cs.repository.GetByID(ctx, id)
	if err != nil {
		log.Printf("error occured while fetching bill with id %d. error %s\n", id, err.Error())
		return &models.Bill{}, err
	}

	return bill, nil
}
