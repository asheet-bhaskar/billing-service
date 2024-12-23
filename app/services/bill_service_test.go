package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
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

func (suite *BillServiceTestSuite) Test_AddLineItemFailsWhenBillNotFound() {
	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      suite.bill.ID,
		Description: "line item 01",
		Amount:      100.0,
		CreatedAt:   time.Now().UTC(),
		Removed:     false,
	}
	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, ce.BillNotFoundError)

	_, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillNotFoundError, err)
}

func (suite *BillServiceTestSuite) Test_AddLineItemFailsWhenBillIsClosed() {
	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      suite.bill.ID,
		Description: "line item 02",
		Amount:      100.0,
		CreatedAt:   time.Now().UTC(),
		Removed:     false,
	}

	bill := *suite.bill
	bill.Status = "closed"

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)

	_, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillClosedError, err)
}

func (suite *BillServiceTestSuite) Test_AddLineItemFailsWhenErrorIsOccurred() {
	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      suite.bill.ID,
		Description: "line item 03",
		Amount:      100.0,
		CreatedAt:   time.Now().UTC(),
		Removed:     false,
	}

	ctx := context.Background()
	testError := errors.New("test-error")
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(suite.bill, nil)
	suite.BillMockRepo.On("AddLineItems", ctx, mock.Anything).Return(lineItem, testError)

	_, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().NotNil(err)
	suite.Require().Equal(testError, err)
}

func (suite *BillServiceTestSuite) Test_AddLineItemSucceeds() {
	lineItem := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      suite.bill.ID,
		Description: "line item 03",
		Amount:      100.0,
		CreatedAt:   time.Now().UTC(),
		Removed:     false,
	}

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(suite.bill, nil)
	suite.BillMockRepo.On("AddLineItems", ctx, mock.Anything).Return(lineItem, nil)

	lineItemSaved, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().Nil(err)
	suite.Require().Equal(lineItem, lineItemSaved)
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
