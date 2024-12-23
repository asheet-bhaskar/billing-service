package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CurrencyRepositoryTestSuite struct {
	suite.Suite
	dbClient *gorm.DB
	cr       CurrencyRepository
}

func (suite *CurrencyRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.dbClient = dbClient.TestDB

	suite.cr = NewCurrencyRepository(dbClient.DB)
}

func (suite *CurrencyRepositoryTestSuite) TearDownSuite() {
	fmt.Printf("cleaning up db records")
	suite.dbClient.Exec("DELETE FROM currencies")
}

func (suite *CurrencyRepositoryTestSuite) Test_CreateCurrencyWhenSucceeds() {
	currency := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "001",
		Name:      "code01",
		Symbol:    "code01",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err := suite.cr.Create(context.Background(), currency)
	suite.Nil(err, "error should be nil")
}

func (suite *CurrencyRepositoryTestSuite) Test_GetCurrencyByIDWhenSucceeds() {
	currency := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "002",
		Name:      "code02",
		Symbol:    "code02",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	currency, err := suite.cr.Create(context.Background(), currency)
	suite.Nil(err, "error should be nil")

	currencyRecord, err := suite.cr.GetByID(context.Background(), currency.ID)

	suite.Nil(err, "error should be nil")
	suite.Equal("002", currencyRecord.Code)
	suite.Equal("code02", currencyRecord.Name)
	suite.Equal("code02", currencyRecord.Symbol)
}

func (suite *CurrencyRepositoryTestSuite) Test_GetCurrencyByCodeWhenSucceeds() {
	currency := &models.Currency{
		ID:        utils.GetNewUUID(),
		Code:      "003",
		Name:      "code03",
		Symbol:    "code03",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	currency, err := suite.cr.Create(context.Background(), currency)
	suite.Nil(err, "error should be nil")

	currencyRecord, err := suite.cr.GetByCode(context.Background(), currency.Code)

	suite.Nil(err, "error should be nil")
	suite.Equal("003", currencyRecord.Code)
	suite.Equal("code03", currencyRecord.Name)
	suite.Equal("code03", currencyRecord.Symbol)
}

func TestCurrencyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyRepositoryTestSuite))
}
