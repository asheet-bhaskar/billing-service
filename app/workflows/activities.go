package workflows

import (
	"context"
	"errors"
	"log"

	"github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/db/repository"
)

type Activities struct {
}

func (a *Activities) AddLineItemActivity(ctx context.Context, message LineItemSignal) error {
	log.Printf("line item %s added, updating the bill amount\n", message.ItemID)

	billRepository := repository.NewBillRepository(db.Clients.DB)
	bill, err := billRepository.GetByID(ctx, message.BillID)

	if err != nil {
		log.Println("error occured while fetching the bill")
		return errors.New("error occured while fetching the bill")
	}

	if bill.Status == "closed" {
		log.Println("already closed bill can not be updated")
		return errors.New("already closed bill can not be updated")
	}

	lineItem, err := billRepository.GetLineItemByID(ctx, message.ItemID)

	if err != nil {
		log.Println("error occured while fetching the bill")
		return errors.New("error occured while fetching the bill")
	}

	updatedAmount := bill.TotalAmount + lineItem.Amount

	err = billRepository.UpdateBillAmount(ctx, message.BillID, updatedAmount)

	if err != nil {
		log.Println("failed to update bill amount")
		return errors.New("failed to update bill amount")
	}

	return nil
}

func (a *Activities) RemoveLineItemActivity(ctx context.Context, message LineItemSignal) error {
	log.Printf("line item removed %s, updating the bill amount\n", message.ItemID)

	billRepository := repository.NewBillRepository(db.Clients.DB)
	bill, err := billRepository.GetByID(ctx, message.BillID)
	if err != nil {
		log.Println("error occured while fetching the bill")
		return errors.New("error occured while fetching the bill")
	}

	if bill.Status == "closed" {
		log.Println("already closed bill can not be updated")
		return errors.New("already closed bill can not be updated")
	}

	lineItem, err := billRepository.GetLineItemByID(ctx, message.ItemID)

	if err != nil {
		log.Println("error occured while fetching the bill")
		return errors.New("error occured while fetching the bill")
	}

	updatedAmount := bill.TotalAmount - lineItem.Amount

	err = billRepository.UpdateBillAmount(ctx, message.BillID, updatedAmount)

	if err != nil {
		log.Println("failed to update bill amount")
		return errors.New("failed to update bill amount")
	}
	return nil
}
