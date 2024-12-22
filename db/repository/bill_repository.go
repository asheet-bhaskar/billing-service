package repository

import (
	"context"
	"log"

	models "github.com/asheet-bhaskar/billing-service/app/services/bill"
	"gorm.io/gorm"
)

type BillRepository struct {
	db *gorm.DB
}

type billRepository interface {
	Create(context.Context, *models.Bill) (*models.Bill, error)
	GetByID(context.Context, int64) (*models.Bill, error)
}

func NewBillRepository(dbClient *gorm.DB) billRepository {
	return &BillRepository{
		db: dbClient,
	}
}

func (br *BillRepository) Create(ctx context.Context, bill *models.Bill) (*models.Bill, error) {
	result := br.db.Create(&bill)

	if result.Error != nil {
		log.Printf("error occured while creating bill, %v. error is %s", bill, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}

func (br *BillRepository) GetByID(ctx context.Context, id int64) (*models.Bill, error) {
	bill := &models.Bill{}
	result := br.db.First(&bill, id)

	if result.Error != nil {
		log.Printf("error occured while querying bill, %d. error is %s", id, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}
