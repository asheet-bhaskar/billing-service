package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	tc "github.com/asheet-bhaskar/billing-service/app/workflows/temporal"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BillServiceTestSuite struct {
	suite.Suite
	BillMockRepo       *repository.MockBillRepository
	CustomerMockRepo   *repository.MockCustomerRepository
	CurrencyMockRepo   *repository.MockCurrencyRepository
	TemporalClientMock *tc.MockTemporalClient
	bs                 BillService
	billRequest        *models.BillRequest
	bill               *models.Bill
	currencyID         string
	customerID         string
}

func (suite *BillServiceTestSuite) SetupTest() {

	billMockRepo := new(repository.MockBillRepository)
	customerMockRepo := new(repository.MockCustomerRepository)
	currencyMockRepo := new(repository.MockCurrencyRepository)
	temporalClientMock := new(tc.MockTemporalClient)

	suite.BillMockRepo = billMockRepo
	suite.CustomerMockRepo = customerMockRepo
	suite.CurrencyMockRepo = currencyMockRepo
	suite.TemporalClientMock = temporalClientMock

	suite.bs = NewBillService(billMockRepo, currencyMockRepo, customerMockRepo, temporalClientMock)
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
	suite.TemporalClientMock.On("ExecuteWorkflow", ctx, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {})

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
	suite.TemporalClientMock.On("SignalWorkflow", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lineItemSaved, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().Nil(err)
	suite.Require().Equal(lineItem, lineItemSaved)
}

func (suite *BillServiceTestSuite) Test_RemoveLineItemFailsWhenBillNotFound() {
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
	suite.BillMockRepo.On("GetLineItemByID", ctx, mock.Anything).Return(lineItem, nil)

	_, err := suite.bs.RemoveLineItems(ctx, "", lineItem.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillNotFoundError, err)
}

func (suite *BillServiceTestSuite) Test_RemoveLineItemFailsWhenBillIsClosed() {
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
	suite.BillMockRepo.On("GetLineItemByID", ctx, mock.Anything).Return(lineItem, nil)

	_, err := suite.bs.RemoveLineItems(ctx, "", lineItem.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillClosedError, err)
}

func (suite *BillServiceTestSuite) Test_RemoveLineItemFailsWhenErrorIsOccurred() {
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
	suite.BillMockRepo.On("GetLineItemByID", ctx, mock.Anything).Return(lineItem, nil)

	_, err := suite.bs.AddLineItems(ctx, lineItem)
	suite.Require().NotNil(err)
	suite.Require().Equal(testError, err)
}

func (suite *BillServiceTestSuite) Test_RemoveLineItemSucceeds() {
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
	suite.BillMockRepo.On("RemoveLineItems", ctx, mock.Anything).Return(lineItem, nil)
	suite.BillMockRepo.On("GetLineItemByID", ctx, mock.Anything).Return(lineItem, nil)
	suite.TemporalClientMock.On("SignalWorkflow", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lineItemSaved, err := suite.bs.RemoveLineItems(ctx, "", lineItem.ID)
	suite.Require().Nil(err)
	suite.Require().Equal(lineItem, lineItemSaved)
}

func (suite *BillServiceTestSuite) Test_CloseBillFailsWhenBillNotFound() {
	bill := *suite.bill

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Bill{}, ce.BillNotFoundError)

	_, err := suite.bs.Close(ctx, bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillNotFoundError, err)
}

func (suite *BillServiceTestSuite) Test_CloseBillFailsWhenBillIsClosed() {
	bill := *suite.bill
	bill.Status = "closed"

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)

	_, err := suite.bs.Close(ctx, bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillClosedError, err)
}

func (suite *BillServiceTestSuite) Test_CloseBillFailsWhenErrorIsOccurred() {
	bill := *suite.bill
	testError := errors.New("test error")
	ctx := context.Background()

	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.BillMockRepo.On("Close", ctx, mock.Anything).Return(&bill, testError)

	_, err := suite.bs.Close(ctx, bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(testError, err)
}

func (suite *BillServiceTestSuite) Test_CloseBillSucceeds() {
	bill := *suite.bill
	closedBill := *suite.bill
	closedBill.Status = "closed"

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.BillMockRepo.On("Close", ctx, mock.Anything).Return(&closedBill, nil)

	billActual, err := suite.bs.Close(ctx, suite.bill.ID)
	suite.Require().Nil(err)
	suite.Require().Equal("closed", billActual.Status)
}

func (suite *BillServiceTestSuite) Test_InvoiceFailsWhenBillNotFound() {
	bill := *suite.bill

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, ce.BillNotFoundError)

	_, err := suite.bs.Invoice(ctx, suite.bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.BillNotFoundError, err)
}

func (suite *BillServiceTestSuite) Test_InvoiceFailsWhenCurrencyNotFound() {
	bill := *suite.bill

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.CurrencyMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Currency{}, ce.CurrencyNotFoundError)

	_, err := suite.bs.Invoice(ctx, suite.bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(ce.CurrencyNotFoundError, err)
}

func (suite *BillServiceTestSuite) Test_InvoiceFailsErrorOccuredWhileFetchingLineItems() {
	bill := *suite.bill
	testError := errors.New("test error")

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.CurrencyMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Currency{Code: "001"}, nil)
	suite.BillMockRepo.On("GetLineItemsByBillID", ctx, mock.Anything).Return([]*models.LineItem{&models.LineItem{}}, testError)

	_, err := suite.bs.Invoice(ctx, suite.bill.ID)
	suite.Require().NotNil(err)
	suite.Require().Equal(testError, err)
}

func (suite *BillServiceTestSuite) Test_InvoiceSucceedsWhenBillNotHasRemovedLineItems() {
	bill := *suite.bill
	lineItems := []*models.LineItem{
		{
			ID:          utils.GetNewUUID(),
			BillID:      bill.ID,
			Description: "line item 001",
			Amount:      100.0,
			CreatedAt:   time.Now(),
			Removed:     false,
		}}

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.CurrencyMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Currency{Code: "001"}, nil)
	suite.BillMockRepo.On("GetLineItemsByBillID", ctx, mock.Anything).Return(lineItems, nil)

	lineItemsActual, err := suite.bs.Invoice(ctx, suite.bill.ID)
	suite.Require().Nil(err)
	suite.Require().Equal(len(lineItems), len(lineItemsActual.LineItems))
}

func (suite *BillServiceTestSuite) Test_InvoiceSucceedsWhenBillHasRemovedLineItems() {
	bill := *suite.bill
	lineItems := []*models.LineItem{
		{
			ID:          utils.GetNewUUID(),
			BillID:      bill.ID,
			Description: "line item 001",
			Amount:      100.0,
			CreatedAt:   time.Now(),
			Removed:     true,
		}}

	ctx := context.Background()
	suite.BillMockRepo.On("GetByID", ctx, mock.Anything).Return(&bill, nil)
	suite.CurrencyMockRepo.On("GetByID", ctx, mock.Anything).Return(&models.Currency{Code: "001"}, nil)
	suite.BillMockRepo.On("GetLineItemsByBillID", ctx, mock.Anything).Return(lineItems, nil)

	lineItemsActual, err := suite.bs.Invoice(ctx, suite.bill.ID)
	suite.Require().Nil(err)
	suite.Require().Equal(0, len(lineItemsActual.LineItems))
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
