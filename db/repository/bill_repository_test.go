package repository

import (
	"context"
	"testing"
	"time"

	billModels "github.com/asheet-bhaskar/billing-service/app/services/bill"
	currencyModels "github.com/asheet-bhaskar/billing-service/app/services/currency"
	customerModels "github.com/asheet-bhaskar/billing-service/app/services/customer"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/stretchr/testify/suite"
)

type BillRepositoryTestSuite struct {
	suite.Suite
	br       billRepository
	csr      customerRepository
	crr      currencyRepository
	bill     *billModels.Bill
	customer *customerModels.Customer
	currency *currencyModels.Currency
}

func (suite *BillRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.br = NewBillRepository(dbClient.DB)
	suite.csr = NewCustomerRepository(dbClient.DB)
	suite.crr = NewCurrencyRepository(dbClient.DB)
	suite.bill = &billModels.Bill{
		Description: "Bill 01",
		CustomerID:  1,
		CurrencyID:  1,
		Status:      "open",
		TotalAmount: 100.00,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().UTC().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	suite.customer = &customerModels.Customer{
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacon@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currency = &currencyModels.Currency{
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type Bill struct {
	ID          int64
	Description string
	CustomerID  int64
	CurrencyID  string
	Status      string
	TotalAmount float64
	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (suite *BillRepositoryTestSuite) Test_CreateBillWhenSucceeds() {
	customer, err := suite.csr.Create(context.Background(), suite.customer)
	suite.Nil(err, "error should be nil")
	suite.NotNil(customer.ID)

	currency, err := suite.crr.Create(context.Background(), suite.currency)
	suite.Nil(err, "error should be nil")
	suite.NotNil(currency.ID)

	suite.bill.CustomerID = customer.ID
	bill, err := suite.br.Create(context.Background(), suite.bill)
	suite.Nil(err, "error should be nil")
	suite.NotNil(bill.ID)
}

// func (suite *BillRepositoryTestSuite) Test_GetCurrencyByIDWhenSucceeds() {
// 	currencyRecord, err := suite.br.GetByID(context.Background(), 1)

// 	suite.Nil(err, "error should be nil")
// 	suite.Equal(int64(1), currencyRecord.ID)
// 	suite.Equal("USD", currencyRecord.Code)
// 	suite.Equal("United states dollar", currencyRecord.Name)
// 	suite.Equal("$", currencyRecord.Symbol)
// }

func TestBillRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BillRepositoryTestSuite))
}
