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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type currencyHandlerTestSuite struct {
	suite.Suite
	billServiceMock     service.BillServiceMock
	customerServiceMock service.CustomerServiceMock
	currencyServiceMock service.CurrencyServiceMock
	apiService          *APIService
}

func (suite *currencyHandlerTestSuite) SetupTest() {
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

func (suite *currencyHandlerTestSuite) Test_CreateReturnSucceeds() {
	ctx := context.Background()
	currencyRequest := &models.Currency{
		Code:   "USD",
		Name:   "United states dollar",
		Symbol: "$",
	}

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currencyServiceMock.On("Create", ctx, mock.Anything).Return(currencyResponse, nil)

	_, err := suite.apiService.CreateCurrencyHandler(ctx, currencyRequest)
	suite.Nil(err)

}

func (suite *currencyHandlerTestSuite) Test_CreateReturnFailsWhenServiceReturnsAlreadyExistError() {
	ctx := context.Background()
	currencyRequest := &models.Currency{
		Code:   "USD",
		Name:   "United states dollar",
		Symbol: "$",
	}

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currencyServiceMock.On("Create", ctx, mock.Anything).Return(currencyResponse, ce.CustomerAlreadyExistError)

	_, err := suite.apiService.CreateCurrencyHandler(ctx, currencyRequest)
	suite.NotNil(err)
}

func (suite *currencyHandlerTestSuite) Test_CreateReturnFailsWhenServiceReturnsUnknownError() {
	ctx := context.Background()
	currencyRequest := &models.Currency{
		Code:   "USD",
		Name:   "United states dollar",
		Symbol: "$",
	}

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	testError := errors.New("test error")

	suite.currencyServiceMock.On("Create", ctx, mock.Anything).Return(currencyResponse, testError)

	_, err := suite.apiService.CreateCurrencyHandler(ctx, currencyRequest)
	suite.NotNil(err)
}

func (suite *currencyHandlerTestSuite) Test_CreateReturnFailsWhenServiceRequestIsInvalid() {
	ctx := context.Background()
	currencyRequest := &models.Currency{
		Code:   "",
		Name:   "",
		Symbol: "$",
	}

	_, err := suite.apiService.CreateCurrencyHandler(ctx, currencyRequest)
	suite.NotNil(err)
}

func (suite *currencyHandlerTestSuite) Test_GetCustoemrHandlerSucceeds() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currencyServiceMock.On("GetByID", ctx, id).Return(currencyResponse, nil)

	_, err := suite.apiService.GetCurrencyHandler(ctx, id)
	suite.Nil(err)
}

func (suite *currencyHandlerTestSuite) Test_GetCustoemrHandlerFailsWhenCustomerNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currencyServiceMock.On("GetByID", ctx, id).Return(currencyResponse, ce.CustomerNotFoundError)

	_, err := suite.apiService.GetCurrencyHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *currencyHandlerTestSuite) Test_GetCustoemrHandlerFailsWhenUnknownErrorOccurs() {
	ctx := context.Background()
	id := utils.GetNewUUID()
	testError := errors.New("test error")

	currencyResponse := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "USD",
		Name:      "United states dollar",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.currencyServiceMock.On("GetByID", ctx, id).Return(currencyResponse, testError)

	_, err := suite.apiService.GetCurrencyHandler(ctx, id)
	suite.NotNil(err)
}

func TestCurrencyHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(currencyHandlerTestSuite))
}
