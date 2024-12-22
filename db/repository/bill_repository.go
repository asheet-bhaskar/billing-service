package repository

import (
	"context"
	"log"

	"github.com/asheet-bhaskar/billing-service/app/models"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"gorm.io/gorm"
)

type billRepository struct {
	db *gorm.DB
}

type BillRepository interface {
	Create(context.Context, *models.Bill) (*models.Bill, error)
	GetByID(context.Context, int64) (*models.Bill, error)
}

func NewBillRepository(dbClient *gorm.DB) BillRepository {
	return &billRepository{
		db: dbClient,
	}
}

func (br *billRepository) Create(ctx context.Context, bill *models.Bill) (*models.Bill, error) {
	result := br.db.Create(&bill)

	if result.Error != nil {
		log.Printf("error occured while creating bill, %v. error is %s", bill, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}

func (br *billRepository) GetByID(ctx context.Context, id int64) (*models.Bill, error) {
	bill := &models.Bill{}
	result := br.db.First(&bill, id)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("bill does not exist for id, %d. error is %s", id, result.Error.Error())
		return bill, ce.BillNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while querying bill, %d. error is %s", id, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}
