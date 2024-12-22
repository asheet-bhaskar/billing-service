package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
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
}

func (suite *BillServiceTestSuite) SetupTest() {

	billMockRepo := new(repository.MockBillRepository)
	customerMockRepo := new(repository.MockCustomerRepository)
	currencyMockRepo := new(repository.MockCurrencyRepository)

	suite.BillMockRepo = billMockRepo
	suite.CustomerMockRepo = customerMockRepo
	suite.CurrencyMockRepo = currencyMockRepo

	suite.bs = NewBillService(billMockRepo, currencyMockRepo, customerMockRepo)

	suite.billRequest = &models.BillRequest{
		Description:  "Bill - 01",
		CustomerID:   int64(100),
		CurrencyCode: "USD",
		PeriodStart:  time.Now().UTC(),
		PeriodEnd:    time.Now().Add(time.Hour * 100),
	}

	suite.bill = &models.Bill{
		Description: "Bill - 01",
		CustomerID:  int64(100),
		CurrencyID:  int64(1),
		Status:      "open",
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().Add(time.Hour * 100),
	}
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.CurrencyMockRepo.On("GetByCode", ctx, "USD").Return(&models.Currency{ID: int64(1)}, nil)
	suite.CustomerMockRepo.On("GetByID", ctx, int64(100)).Return(&models.Customer{ID: int64(100)}, nil)
	suite.BillMockRepo.On("Create", ctx, mock.Anything).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.Create(ctx, suite.billRequest)

	suite.Require().Error(err)
	suite.Require().NotEqual(0, bill.ID)
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.CustomerMockRepo.On("GetByID", ctx, int64(100)).Return(&models.Customer{}, nil)
	suite.CurrencyMockRepo.On("GetByCode", ctx, "USD").Return(&models.Currency{ID: 1}, nil)
	suite.BillMockRepo.On("Create", ctx, mock.Anything).Return(&models.Bill{}, nil)

	bill, err := suite.bs.Create(ctx, suite.billRequest)

	suite.Require().Nil(err)
	suite.Require().Equal(&models.Bill{}, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsErrorWhenFails() {
	ctx := context.Background()
	billID := int64(1)
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.GetByID(ctx, billID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.bill, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	billID := int64(1)
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, nil)

	bill, err := suite.bs.GetByID(ctx, billID)

	suite.Require().Nil(err)
	suite.Require().Equal(&models.Bill{}, bill)
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
