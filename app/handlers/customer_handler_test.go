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

type customerHandlerTestSuite struct {
	suite.Suite
	billServiceMock     service.BillServiceMock
	customerServiceMock service.CustomerServiceMock
	currencyServiceMock service.CurrencyServiceMock
	apiService          *APIService
}

func (suite *customerHandlerTestSuite) SetupTest() {
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

func (suite *customerHandlerTestSuite) Test_CreateReturnSucceeds() {
	ctx := context.Background()
	customerRequest := &models.Customer{
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
	}

	customerResponse := &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.customerServiceMock.On("Create", ctx, mock.Anything).Return(customerResponse, nil)

	customer, err := suite.apiService.CreateCustomerHandler(ctx, customerRequest)
	suite.Nil(err)
	suite.Equal(customerResponse.ID, customer.ID)
	suite.Equal(customerResponse.CreatedAt, customer.CreatedAt)
	suite.Equal(customerResponse.UpdatedAt, customer.UpdatedAt)
	suite.Equal(customerResponse.FirstName, customer.FirstName)
	suite.Equal(customerResponse.LastName, customer.LastName)
	suite.Equal(customerResponse.Email, customer.Email)
}

func (suite *customerHandlerTestSuite) Test_CreateReturnFailsWhenServiceReturnsAlreadyExistError() {
	ctx := context.Background()
	customerRequest := &models.Customer{
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
	}

	customerResponse := &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.customerServiceMock.On("Create", ctx, mock.Anything).Return(customerResponse, ce.CustomerAlreadyExistError)

	_, err := suite.apiService.CreateCustomerHandler(ctx, customerRequest)
	suite.NotNil(err)
}

func (suite *customerHandlerTestSuite) Test_CreateReturnFailsWhenServiceReturnsUnknownError() {
	ctx := context.Background()
	customerRequest := &models.Customer{
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
	}

	customerResponse := &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	testError := errors.New("test error")

	suite.customerServiceMock.On("Create", ctx, mock.Anything).Return(customerResponse, testError)

	_, err := suite.apiService.CreateCustomerHandler(ctx, customerRequest)
	suite.NotNil(err)
}

func (suite *customerHandlerTestSuite) Test_CreateReturnFailsWhenServiceRequestIsInvalid() {
	ctx := context.Background()
	customerRequest := &models.Customer{
		FirstName: "",
		LastName:  "",
		Email:     "john.jacobs@mail.com",
	}

	_, err := suite.apiService.CreateCustomerHandler(ctx, customerRequest)
	suite.NotNil(err)
}

func (suite *customerHandlerTestSuite) Test_GetCustoemrHandlerSucceeds() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	customerResponse := &models.Customer{
		ID:        id,
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.customerServiceMock.On("GetByID", ctx, id).Return(customerResponse, nil)

	_, err := suite.apiService.GetCustomerHandler(ctx, id)
	suite.Nil(err)
}

func (suite *customerHandlerTestSuite) Test_GetCustoemrHandlerFailsWhenCustomerNotFound() {
	ctx := context.Background()
	id := utils.GetNewUUID()

	customerResponse := &models.Customer{
		ID:        id,
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.customerServiceMock.On("GetByID", ctx, id).Return(customerResponse, ce.CustomerNotFoundError)

	_, err := suite.apiService.GetCustomerHandler(ctx, id)
	suite.NotNil(err)
}

func (suite *customerHandlerTestSuite) Test_GetCustoemrHandlerFailsWhenUnknownErrorOccurs() {
	ctx := context.Background()
	id := utils.GetNewUUID()
	testError := errors.New("test error")

	customerResponse := &models.Customer{
		ID:        id,
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacobs@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	suite.customerServiceMock.On("GetByID", ctx, id).Return(customerResponse, testError)

	_, err := suite.apiService.GetCustomerHandler(ctx, id)
	suite.NotNil(err)
}

func TestCustomerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(customerHandlerTestSuite))
}
