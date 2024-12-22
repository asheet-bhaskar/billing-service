package repository

import (
	"context"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/stretchr/testify/suite"
)

type CurrencyRepositoryTestSuite struct {
	suite.Suite
	cr       CurrencyRepository
	currency *models.Currency
}

func (suite *CurrencyRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.cr = NewCurrencyRepository(dbClient.DB)
	suite.currency = &models.Currency{
		Code:      "GEL",
		Name:      "Geogrian Lari",
		Symbol:    "ლ",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (suite *CurrencyRepositoryTestSuite) Test_CreateCurrencyWhenSucceeds() {
	currency, err := suite.cr.Create(context.Background(), suite.currency)
	suite.Nil(err, "error should be nil")
	suite.NotNil(currency.ID)
}

func (suite *CurrencyRepositoryTestSuite) Test_GetCurrencyByIDWhenSucceeds() {
	currencyRecord, err := suite.cr.GetByID(context.Background(), 1)

	suite.Nil(err, "error should be nil")
	suite.Equal(int64(1), currencyRecord.ID)
	suite.Equal("USD", currencyRecord.Code)
	suite.Equal("United states dollar", currencyRecord.Name)
	suite.Equal("$", currencyRecord.Symbol)
}

func (suite *CurrencyRepositoryTestSuite) Test_GetCurrencyByIDWhenFails() {
	currencyRecord, err := suite.cr.GetByCode(context.Background(), "USD")

	suite.Nil(err, "error should be nil")
	suite.Equal(int64(1), currencyRecord.ID)
	suite.Equal("USD", currencyRecord.Code)
	suite.Equal("United states dollar", currencyRecord.Name)
	suite.Equal("$", currencyRecord.Symbol)
}

func TestCurrencyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyRepositoryTestSuite))
}
