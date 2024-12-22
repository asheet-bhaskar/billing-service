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

type MockCurrencyRepository struct {
	mock.Mock
}

func (m *MockCurrencyRepository) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) GetByID(ctx context.Context, id int64) (*models.Currency, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) GetByCode(ctx context.Context, id string) (*models.Currency, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Currency), args.Error(1)
}

type CurrencyServiceTestSuite struct {
	suite.Suite
	MockRepo *MockCurrencyRepository
	cs       CurrencyService
	currency *models.Currency
}

func (suite *CurrencyServiceTestSuite) SetupTest() {

	mockRepo := new(MockCurrencyRepository)
	suite.MockRepo = mockRepo
	suite.cs = NewCurrencyService(mockRepo)

	suite.currency = &models.Currency{
		Code:      "USD",
		Name:      "United states dollars",
		Symbol:    "$",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (suite *CurrencyServiceTestSuite) Test_CreateCurrencyReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.currency).Return(&models.Currency{}, errors.New("test-error"))

	currency, err := suite.cs.Create(ctx, suite.currency)

	suite.Require().Error(err)
	suite.Require().NotEqual(suite.currency, currency)
}

func (suite *CurrencyServiceTestSuite) Test_CreateCurrencyReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.currency).Return(suite.currency, nil)

	currency, err := suite.cs.Create(ctx, suite.currency)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.currency, currency)
}

func (suite *CurrencyServiceTestSuite) Test_GetByIDReturnsErrorWhenFails() {
	ctx := context.Background()
	currencyID := int64(1)
	suite.MockRepo.On("GetByID", ctx, currencyID).Return(&models.Currency{}, errors.New("test-error"))

	currency, err := suite.cs.GetByID(ctx, currencyID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.currency, currency)
}

func (suite *CurrencyServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	currencyID := int64(1)
	suite.MockRepo.On("GetByID", ctx, currencyID).Return(suite.currency, nil)

	currency, err := suite.cs.GetByID(ctx, currencyID)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.currency, currency)
}

func TestCurrencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}
