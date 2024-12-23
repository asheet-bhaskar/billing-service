package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BillServiceTestSuite struct {
	suite.Suite
	BillMockRepo     *repository.MockBillRepository
	CustomerMockRepo *repository.MockCustomerRepository
	CurrencyMockRepo *repository.MockCurrencyRepository
	bs               BillService
	billRequest      *models.BillRequest
	bill             *models.Bill
	currencyID       string
	customerID       string
}

func (suite *BillServiceTestSuite) SetupTest() {

	billMockRepo := new(repository.MockBillRepository)
	customerMockRepo := new(repository.MockCustomerRepository)
	currencyMockRepo := new(repository.MockCurrencyRepository)

	suite.BillMockRepo = billMockRepo
	suite.CustomerMockRepo = customerMockRepo
	suite.CurrencyMockRepo = currencyMockRepo

	suite.bs = NewBillService(billMockRepo, currencyMockRepo, customerMockRepo)
	currencyID := utils.GetNewUUID()
	customerID := utils.GetNewUUID()

	suite.billRequest = &models.BillRequest{
		Description:  "Bill - 01",
		CustomerID:   customerID,
		CurrencyCode: "USD",
		PeriodStart:  time.Now().UTC(),
		PeriodEnd:    time.Now().Add(time.Hour * 100),
	}

	suite.bill = &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "Bill - 01",
		CustomerID:  customerID,
		CurrencyID:  currencyID,
		Status:      "open",
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().Add(time.Hour * 100),
	}

	suite.customerID = customerID
	suite.currencyID = currencyID
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.CurrencyMockRepo.On("GetByCode", ctx, "USD").Return(&models.Currency{ID: suite.currencyID}, nil)
	suite.CustomerMockRepo.On("GetByID", ctx, suite.customerID).Return(&models.Customer{ID: suite.customerID}, nil)
	suite.BillMockRepo.On("Create", ctx, mock.Anything).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.Create(ctx, suite.billRequest)

	suite.Require().Error(err)
	suite.Require().NotEqual(0, bill.ID)
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.CustomerMockRepo.On("GetByID", ctx, suite.customerID).Return(&models.Customer{}, nil)
	suite.CurrencyMockRepo.On("GetByCode", ctx, "USD").Return(&models.Currency{ID: suite.currencyID}, nil)
	suite.BillMockRepo.On("Create", ctx, mock.Anything).Return(&models.Bill{}, nil)

	bill, err := suite.bs.Create(ctx, suite.billRequest)

	suite.Require().Nil(err)
	suite.Require().Equal(&models.Bill{}, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.GetByID(ctx, suite.bill.ID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.bill, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, nil)

	bill, err := suite.bs.GetByID(ctx, suite.bill.ID)

	suite.Require().Nil(err)
	suite.Require().Equal(&models.Bill{}, bill)
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
