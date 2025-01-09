package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	service "github.com/asheet-bhaskar/billing-service/app/services"
	ce "github.com/asheet-bhaskar/billing-service/pkg/error"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type billHandlerTestSuite struct {
	suite.Suite
	billServiceMock     service.BillServiceMock
	customerServiceMock service.CustomerServiceMock
	currencyServiceMock service.CurrencyServiceMock
	apiService          *APIService
}

func (suite *billHandlerTestSuite) SetupTest() {
	billServiceMock := new(service.BillServiceMock)
	customerServiceMock := new(service.CustomerServiceMock)
	currencyServiceMock := new(service.CurrencyServiceMock)

	suite.billServiceMock = *billServiceMock
	suite.currencyServiceMock = *currencyServiceMock
	suite.customerServiceMock = *customerServiceMock
	suite.apiService = &APIService{
		Bill:     &suite.billServiceMock,
		Customer: &suite.customerServiceMock,
		Currency: &suite.currencyServiceMock,
	}
}

func (suite *billHandlerTestSuite) Test_CreateBillHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	billRequest := &models.BillRequest{
		Description:  "bill 01",
		CustomerID:   "customer-id",
		CurrencyCode: "USD",
		PeriodStart:  now,
		PeriodEnd:    now.Add(time.Hour * 24),
	}

	billResponse := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "bill 01",
		CustomerID:  utils.GetNewUUID(),
		CurrencyID:  utils.GetNewUUID(),
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour * 24),
		TotalAmount: 0.0,
		Status:      "open",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	suite.billServiceMock.On("Create", ctx, billRequest).Return(billResponse, nil)

	_, err := suite.apiService.CreateBillHandler(ctx, billRequest)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_CreateBillHandlerFailsWhenUnknownErrorOccurs() {
	ctx := context.Background()
	now := time.Now().UTC()
	billRequest := &models.BillRequest{
		Description:  "bill 01",
		CustomerID:   "customer-id",
		CurrencyCode: "USD",
		PeriodStart:  now,
		PeriodEnd:    now.Add(time.Hour * 24),
	}

	billResponse := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "bill 01",
		CustomerID:  utils.GetNewUUID(),
		CurrencyID:  utils.GetNewUUID(),
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour * 24),
		TotalAmount: 0.0,
		Status:      "open",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	testError := errors.New("test error")

	suite.billServiceMock.On("Create", ctx, billRequest).Return(billResponse, testError)

	_, err := suite.apiService.CreateBillHandler(ctx, billRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CreateBillHandlerFailsWhenCurrencyNotFound() {
	ctx := context.Background()
	now := time.Now().UTC()
	billRequest := &models.BillRequest{
		Description:  "bill 01",
		CustomerID:   "customer-id",
		CurrencyCode: "USD",
		PeriodStart:  now,
		PeriodEnd:    now.Add(time.Hour * 24),
	}

	suite.billServiceMock.On("Create", ctx, billRequest).Return(&models.Bill{}, ce.CurrencyNotFoundError)

	_, err := suite.apiService.CreateBillHandler(ctx, billRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CreateBillHandlerFailsWhenCustomerNotFound() {
	ctx := context.Background()
	now := time.Now().UTC()
	billRequest := &models.BillRequest{
		Description:  "bill 01",
		CustomerID:   "customer-id",
		CurrencyCode: "USD",
		PeriodStart:  now,
		PeriodEnd:    now.Add(time.Hour * 24),
	}

	suite.billServiceMock.On("Create", ctx, billRequest).Return(&models.Bill{}, ce.CustomerNotFoundError)

	_, err := suite.apiService.CreateBillHandler(ctx, billRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CreateBillHandlerFailsWhenRequestIsInvalid() {
	ctx := context.Background()
	now := time.Now().UTC()
	billRequest := &models.BillRequest{
		Description:  "bill 01",
		CustomerID:   "",
		CurrencyCode: "",
		PeriodStart:  now,
		PeriodEnd:    now.Add(time.Hour * 24),
	}

	_, err := suite.apiService.CreateBillHandler(ctx, billRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetBillHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	billResponse := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "bill 01",
		CustomerID:  utils.GetNewUUID(),
		CurrencyID:  utils.GetNewUUID(),
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour * 24),
		TotalAmount: 0.0,
		Status:      "open",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	suite.billServiceMock.On("GetByID", ctx, billResponse.ID).Return(billResponse, nil)

	_, err := suite.apiService.GetBillHandler(ctx, billResponse.ID)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_GetBillHandlerFailsWhenIDIsInvalid() {
	ctx := context.Background()

	_, err := suite.apiService.GetBillHandler(ctx, "")
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetBillHandlerFailsWhenBillIsNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	suite.billServiceMock.On("GetByID", ctx, id).Return(&models.Bill{}, ce.BillNotFoundError)

	_, err := suite.apiService.GetBillHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetBillHandlerFailsWhenUnknownErrorOccured() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	testError := errors.New("test error")
	suite.billServiceMock.On("GetByID", ctx, id).Return(&models.Bill{}, testError)

	_, err := suite.apiService.GetBillHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CloseBillHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	billResponse := &models.Bill{
		ID:          utils.GetNewUUID(),
		Description: "bill 01",
		CustomerID:  utils.GetNewUUID(),
		CurrencyID:  utils.GetNewUUID(),
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour * 24),
		TotalAmount: 0.0,
		Status:      "closed",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	suite.billServiceMock.On("Close", ctx, billResponse.ID).Return(billResponse, nil)

	_, err := suite.apiService.CloseBillHandler(ctx, billResponse.ID)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_CloseBillHandlerFailsWhenIDIsInvalid() {
	ctx := context.Background()

	_, err := suite.apiService.CloseBillHandler(ctx, "")
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CloseBillHandlerFailsWhenBillIsNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	suite.billServiceMock.On("Close", ctx, id).Return(&models.Bill{}, ce.BillNotFoundError)

	_, err := suite.apiService.CloseBillHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_CloseBillHandlerFailsWhenUnknownErrorOccured() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	testError := errors.New("test error")
	suite.billServiceMock.On("Close", ctx, id).Return(&models.Bill{}, testError)

	_, err := suite.apiService.CloseBillHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetInvoiceHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	id := utils.GetNewUUID()
	invoice := &models.Invoice{
		BillID:      id,
		Description: "bill 01",
		CustomerID:  utils.GetNewUUID(),
		CurrencyID:  utils.GetNewUUID(),
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour * 24),
		TotalAmount: 0.0,
		Status:      "closed",
		LineItems:   []models.LineItem{},
	}

	suite.billServiceMock.On("Invoice", ctx, id).Return(invoice, nil)

	_, err := suite.apiService.GetInvoiceHandler(ctx, id)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_GetInvoiceHandlerFailsWhenIDIsInvalid() {
	ctx := context.Background()

	_, err := suite.apiService.GetInvoiceHandler(ctx, "")
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetInvoiceHandlerFailsWhenBillIsNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	suite.billServiceMock.On("Invoice", ctx, id).Return(&models.Invoice{}, ce.BillNotFoundError)

	_, err := suite.apiService.GetInvoiceHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_GetInvoiceHandlerFailsWhenUnknownErrorOccured() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	testError := errors.New("test error")
	suite.billServiceMock.On("Invoice", ctx, id).Return(&models.Invoice{}, testError)

	_, err := suite.apiService.GetInvoiceHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_AddLineItemHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	id := utils.GetNewUUID()
	lineItemRequest := &models.AddLineItemrequest{
		BillID:      id,
		Description: "item 01",
		Amount:      10.0,
	}

	lineItemResponse := &models.LineItem{
		ID:          utils.GetNewUUID(),
		BillID:      id,
		Description: "item 01",
		Amount:      10.0,
		CreatedAt:   now,
		Removed:     false,
	}

	suite.billServiceMock.On("AddLineItems", ctx, lineItemRequest.ToLineItem()).Return(lineItemResponse, nil)

	_, err := suite.apiService.AddLineItemsHandler(ctx, *lineItemRequest)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_AddLineItemHandlerFailsWhenLineItemIsInvalid() {
	ctx := context.Background()
	lineItemRequest := &models.AddLineItemrequest{
		BillID:      "",
		Description: "item 01",
		Amount:      10.0,
	}

	_, err := suite.apiService.AddLineItemsHandler(ctx, *lineItemRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_AddLineItemHandlerFailWhenBillNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()
	lineItemRequest := &models.AddLineItemrequest{
		BillID:      id,
		Description: "item 01",
		Amount:      10.0,
	}

	suite.billServiceMock.On("AddLineItems", ctx, lineItemRequest.ToLineItem()).Return(&models.LineItem{}, ce.BillNotFoundError)

	_, err := suite.apiService.AddLineItemsHandler(ctx, *lineItemRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_AddLineItemHandlerFailWhenBillIsClosed() {
	ctx := context.Background()
	id := utils.GetNewUUID()
	lineItemRequest := &models.AddLineItemrequest{
		BillID:      id,
		Description: "item 01",
		Amount:      10.0,
	}

	suite.billServiceMock.On("AddLineItems", ctx, lineItemRequest.ToLineItem()).Return(&models.LineItem{}, ce.BillClosedError)

	_, err := suite.apiService.AddLineItemsHandler(ctx, *lineItemRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_AddLineItemHandlerFailWhenUnknownErrorOccurs() {
	ctx := context.Background()
	id := utils.GetNewUUID()
	lineItemRequest := &models.AddLineItemrequest{
		BillID:      id,
		Description: "item 01",
		Amount:      10.0,
	}

	testError := errors.New("test error")

	suite.billServiceMock.On("AddLineItems", ctx, lineItemRequest.ToLineItem()).Return(&models.LineItem{}, testError)

	_, err := suite.apiService.AddLineItemsHandler(ctx, *lineItemRequest)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_RemoveLineItemHandlerSucceeds() {
	ctx := context.Background()
	now := time.Now().UTC()
	itemID := utils.GetNewUUID()
	billID := utils.GetNewUUID()

	lineItemResponse := &models.LineItem{
		ID:          itemID,
		BillID:      billID,
		Description: "item 01",
		Amount:      10.0,
		CreatedAt:   now,
		Removed:     false,
	}

	suite.billServiceMock.On("RemoveLineItems", ctx, billID, itemID).Return(lineItemResponse, nil)

	_, err := suite.apiService.RemoveLineItemsHandler(ctx, billID, itemID)
	suite.Nil(err)
}

func (suite *billHandlerTestSuite) Test_RemoveLineItemHandlerFailsWhenitemIDIsInvalid() {
	ctx := context.Background()
	itemID := ""
	billID := utils.GetNewUUID()

	_, err := suite.apiService.RemoveLineItemsHandler(ctx, billID, itemID)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_RemoveLineItemHandlerFailWhenBillNotFound() {
	ctx := context.Background()
	itemID := utils.GetNewUUID()
	billID := utils.GetNewUUID()

	suite.billServiceMock.On("RemoveLineItems", ctx, billID, itemID).Return(&models.LineItem{}, ce.BillNotFoundError)

	_, err := suite.apiService.RemoveLineItemsHandler(ctx, billID, itemID)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_RemoveLineItemHandlerFailWhenBillIsClosed() {
	ctx := context.Background()
	itemID := utils.GetNewUUID()
	billID := utils.GetNewUUID()

	suite.billServiceMock.On("RemoveLineItems", ctx, billID, itemID).Return(&models.LineItem{}, ce.BillClosedError)

	_, err := suite.apiService.RemoveLineItemsHandler(ctx, billID, itemID)
	suite.NotNil(err)
}

func (suite *billHandlerTestSuite) Test_RemoveLineItemHandlerFailWhenUnknownErrorOccurs() {
	ctx := context.Background()
	itemID := utils.GetNewUUID()
	billID := utils.GetNewUUID()

	testError := errors.New("test error")

	suite.billServiceMock.On("RemoveLineItems", ctx, billID, itemID).Return(&models.LineItem{}, testError)

	_, err := suite.apiService.RemoveLineItemsHandler(ctx, billID, itemID)
	suite.NotNil(err)
}

func TestBillHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(billHandlerTestSuite))
}
