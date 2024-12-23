package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/asheet-bhaskar/billing-service/db/repository"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type CustomerServiceTestSuite struct {
	suite.Suite
	MockRepo *repository.MockCustomerRepository
	cs       CustomerService
	customer *models.Customer
}

func (suite *CustomerServiceTestSuite) SetupTest() {

	mockRepo := new(repository.MockCustomerRepository)
	suite.MockRepo = mockRepo
	suite.cs = NewCustomerService(mockRepo)

	suite.customer = &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacon@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (suite *CustomerServiceTestSuite) Test_CreateCustomerReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.customer).Return(&models.Customer{}, errors.New("test-error"))

	customer, err := suite.cs.Create(ctx, suite.customer)

	suite.Require().Error(err)
	suite.Require().NotEqual(suite.customer, customer)
}

func (suite *CustomerServiceTestSuite) Test_CreateCustomerReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.customer).Return(suite.customer, nil)

	customer, err := suite.cs.Create(ctx, suite.customer)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.customer, customer)
}

func (suite *CustomerServiceTestSuite) Test_GetByIDReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.MockRepo.On("GetByID", ctx, suite.customer.ID).Return(&models.Customer{}, errors.New("test-error"))

	customer, err := suite.cs.GetByID(ctx, suite.customer.ID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.customer, customer)
}

func (suite *CustomerServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.MockRepo.On("GetByID", ctx, suite.customer.ID).Return(suite.customer, nil)

	customer, err := suite.cs.GetByID(ctx, suite.customer.ID)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.customer, customer)
}

func TestCustomerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerServiceTestSuite))
}
