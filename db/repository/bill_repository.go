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
	GetByID(context.Context, string) (*models.Bill, error)
	AddLineItems(context.Context, *models.LineItem) (*models.LineItem, error)
	RemoveLineItems(context.Context, *models.LineItem) (*models.LineItem, error)
	GetLineItemsByBillID(context.Context, string) ([]*models.LineItem, error)
	GetLineItemByID(context.Context, string) (*models.LineItem, error)
	Close(context.Context, string) (*models.Bill, error)
	UpdateBillAmount(context.Context, string, float64) error
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

func (br *billRepository) GetByID(ctx context.Context, id string) (*models.Bill, error) {
	bill := &models.Bill{}
	result := br.db.Where("id = ?", id).Find(&bill)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("bill does not exist for id, %s. error is %s", id, result.Error.Error())
		return bill, ce.BillNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while querying bill, %s. error is %s", id, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}

func (br *billRepository) AddLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	result := br.db.Create(&lineItem)

	if result.Error != nil {
		log.Printf("error occured while creating lineItem, %v. error is %s", lineItem, result.Error.Error())
		return lineItem, result.Error
	}

	return lineItem, nil
}

func (br *billRepository) RemoveLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	lineItem.Removed = true
	log.Printf("removing line item %v\n", lineItem)
	result := br.db.Model(&lineItem).Where("id = ?", lineItem.ID).Update("removed", true)

	if result.Error != nil {
		log.Printf("error occured while removing lineItem, %v. error is %s", lineItem, result.Error.Error())
		return lineItem, result.Error
	}

	return lineItem, nil
}

func (br *billRepository) Close(ctx context.Context, id string) (*models.Bill, error) {
	bill, err := br.GetByID(ctx, id)

	if err != nil {
		log.Printf("error occured while fetching bill id, %s. error is %s", id, err.Error())
		return bill, err
	}

	bill.Status = "closed"
	result := br.db.Save(bill)

	if result.Error != nil {
		log.Printf("error occured while closing the bill id, %s. error is %s", id, result.Error.Error())
		return bill, result.Error
	}

	return bill, nil
}

func (br *billRepository) GetLineItemsByBillID(ctx context.Context, billID string) ([]*models.LineItem, error) {
	lineItems := []*models.LineItem{}
	result := br.db.Where("bill_id = ?", billID).Find(&lineItems)

	if result.Error != nil {
		log.Printf("error occured while fetching line items for bill id, %s. error is %s", billID, result.Error.Error())
		return lineItems, result.Error
	}

	return lineItems, nil
}

func (br *billRepository) GetLineItemByID(ctx context.Context, id string) (*models.LineItem, error) {
	lineItem := &models.LineItem{}
	result := br.db.Where("id = ?", id).Find(&lineItem)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("line item does not exist for id, %s. error is %s", id, result.Error.Error())
		return lineItem, ce.LineItemNotFoundError
	}

	if result.Error != nil {
		log.Printf("error occured while fetching line item for id, %s. error is %s", id, result.Error.Error())
		return lineItem, result.Error
	}

	return lineItem, nil
}

func (br *billRepository) UpdateBillAmount(ctx context.Context, billID string, amount float64) error {
	log.Printf("updating bill amount %f for bill id %s\n", amount, billID)
	bill := &models.Bill{}
	result := br.db.Model(bill).Where("id = ?", billID).Update("total_amount", amount)

	if result.Error != nil {
		log.Printf("error occured while updating amount for bill %s\n", billID)
		return result.Error
	}

	return nil
}
