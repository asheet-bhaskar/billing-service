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

type CurrencyServiceTestSuite struct {
	suite.Suite
	MockRepo *repository.MockCurrencyRepository
	cs       CurrencyService
	currency *models.Currency
}

func (suite *CurrencyServiceTestSuite) SetupTest() {

	mockRepo := new(repository.MockCurrencyRepository)
	suite.MockRepo = mockRepo
	suite.cs = NewCurrencyService(mockRepo)

	suite.currency = &models.Currency{
		ID:        utils.GetNewUUID(),
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

	suite.MockRepo.On("GetByID", ctx, suite.currency.ID).Return(&models.Currency{}, errors.New("test-error"))

	currency, err := suite.cs.GetByID(ctx, suite.currency.ID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.currency, currency)
}

func (suite *CurrencyServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()

	suite.MockRepo.On("GetByID", ctx, suite.currency.ID).Return(suite.currency, nil)

	currency, err := suite.cs.GetByID(ctx, suite.currency.ID)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.currency, currency)
}

func TestCurrencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}
