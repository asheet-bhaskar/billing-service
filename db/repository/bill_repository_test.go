package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BillRepositoryTestSuite struct {
	suite.Suite
	dbClient *gorm.DB
	br       BillRepository
	csr      CustomerRepository
	crr      CurrencyRepository
	bill     *models.Bill
	customer *models.Customer
	currency *models.Currency
}

func (suite *BillRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.dbClient = dbClient.TestDB
	suite.br = NewBillRepository(dbClient.DB)
	suite.csr = NewCustomerRepository(dbClient.DB)
	suite.crr = NewCurrencyRepository(dbClient.DB)

	customer := &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacon@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	currency := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      utils.RandomString(3),
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = suite.csr.Create(context.Background(), customer)
	suite.Nil(err, "error should be nil")

	_, err = suite.crr.Create(context.Background(), currency)
	suite.Nil(err, "error should be nil")

	suite.currency = currency
	suite.customer = customer

	suite.bill = &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  customer.ID,
		CurrencyID:  currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (suite *BillRepositoryTestSuite) TearDownSuite() {
	fmt.Printf("cleaning up db records")
	suite.dbClient.Exec("DELETE FROM bills")
	suite.dbClient.Exec("DELETE FROM line_items")
	suite.dbClient.Exec("DELETE FROM currencies")
	suite.dbClient.Exec("DELETE FROM customers")
}

func (suite *BillRepositoryTestSuite) Test_CreateBillWhenSucceeds() {
	_, err := suite.br.Create(context.Background(), suite.bill)
	suite.Nil(err, "error should be nil")
}

func (suite *BillRepositoryTestSuite) Test_GetByIDWhenSucceeds() {
	_, err := suite.br.GetByID(context.Background(), suite.bill.ID)
	suite.Nil(err, "error should be nil")
}

func (suite *BillRepositoryTestSuite) Test_AddLineItemWhenSucceeds() {
	ctx := context.Background()
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  suite.customer.ID,
		CurrencyID:  suite.currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := suite.br.Create(ctx, bill)
	suite.Nil(err, "error should be nil")

	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      bill.ID,
		Description: "line item 01",
		Amount:      12.50,
		CreatedAt:   time.Now(),
		Removed:     false,
	}

	_, err = suite.br.AddLineItems(ctx, lineItem)
	suite.Nil(err, "error should be nil")
}

func (suite *BillRepositoryTestSuite) Test_RemoveLineItemWhenSucceeds() {
	ctx := context.Background()
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  suite.customer.ID,
		CurrencyID:  suite.currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := suite.br.Create(ctx, bill)
	suite.Nil(err, "error should be nil")

	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      bill.ID,
		Description: "line item 01",
		Amount:      12.50,
		CreatedAt:   time.Now(),
		Removed:     false,
	}

	_, err = suite.br.AddLineItems(ctx, lineItem)
	suite.Nil(err, "error should be nil")

	_, err = suite.br.RemoveLineItems(ctx, lineItem)
	suite.Nil(err, "error should be nil")
}

func (suite *BillRepositoryTestSuite) Test_GetLineItemByIDSucceeds() {
	ctx := context.Background()
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  suite.customer.ID,
		CurrencyID:  suite.currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := suite.br.Create(ctx, bill)
	suite.Nil(err, "error should be nil")

	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      bill.ID,
		Description: "line item 01",
		Amount:      12.50,
		CreatedAt:   time.Now(),
		Removed:     false,
	}

	_, err = suite.br.AddLineItems(ctx, lineItem)
	suite.Nil(err, "error should be nil")

	lineItemActual, err := suite.br.GetLineItemByID(ctx, lineItem.ID)
	suite.Nil(err, "error should be nil")
	suite.Equal(lineItem.ID, lineItemActual.ID)
}

func (suite *BillRepositoryTestSuite) Test_GetLineItemByBillIDWhenSucceeds() {
	ctx := context.Background()
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  suite.customer.ID,
		CurrencyID:  suite.currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := suite.br.Create(ctx, bill)
	suite.Nil(err, "error should be nil")

	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      bill.ID,
		Description: "line item 01",
		Amount:      12.50,
		CreatedAt:   time.Now(),
		Removed:     false,
	}

	_, err = suite.br.AddLineItems(ctx, lineItem)
	suite.Nil(err, "error should be nil")

	lineItems, err := suite.br.GetLineItemsByBillID(ctx, bill.ID)
	suite.Nil(err, "error should be nil")
	suite.Equal(1, len(lineItems))
}

func (suite *BillRepositoryTestSuite) Test_CloseDWhenSucceeds() {
	ctx := context.Background()
	bill := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill 01",
		CustomerID:  suite.customer.ID,
		CurrencyID:  suite.currency.ID,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err := suite.br.Create(ctx, bill)
	suite.Nil(err, "error should be nil")

	closeBill, err := suite.br.Close(ctx, bill.ID)
	suite.Nil(err, "error should be nil")
	suite.Equal("closed", closeBill.Status)
}

func TestBillRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BillRepositoryTestSuite))
}
