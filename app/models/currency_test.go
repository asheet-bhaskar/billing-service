package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CurrencyTestSuite struct {
	suite.Suite
	validCurrency   *Currency
	invalidCurrency *Currency
}

func (suite *CurrencyTestSuite) SetupTest() {
	suite.validCurrency = &Currency{
		Code:   "USD",
		Name:   "United states dolalr",
		Symbol: "$",
	}

	suite.invalidCurrency = &Currency{
		Code: "USD",
	}
}

func (suite *CurrencyTestSuite) Test_IsValidReturnFalse() {
	suite.False(suite.invalidCurrency.IsValid())
}

func (suite *CurrencyTestSuite) Test_IsValidReturnTrue() {
	suite.True(suite.validCurrency.IsValid())
}

func TestCurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyTestSuite))
}
