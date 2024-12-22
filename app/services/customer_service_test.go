package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Customer), args.Error(1)
}

type CustomerServiceTestSuite struct {
	suite.Suite
	MockRepo *MockCustomerRepository
	cs       CustomerService
	customer *models.Customer
}

func (suite *CustomerServiceTestSuite) SetupTest() {

	mockRepo := new(MockCustomerRepository)
	suite.MockRepo = mockRepo
	suite.cs = NewCustomerService(mockRepo)

	suite.customer = &models.Customer{
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
	customerID := int64(1)
	suite.MockRepo.On("GetByID", ctx, customerID).Return(&models.Customer{}, errors.New("test-error"))

	customer, err := suite.cs.GetByID(ctx, customerID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.customer, customer)
}

func (suite *CustomerServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	customerID := int64(1)
	suite.MockRepo.On("GetByID", ctx, customerID).Return(suite.customer, nil)

	customer, err := suite.cs.GetByID(ctx, customerID)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.customer, customer)
}

func TestCustomerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerServiceTestSuite))
}
