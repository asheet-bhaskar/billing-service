package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type BillTestSuite struct {
	suite.Suite
	validBillrequest   *BillRequest
	invalidBillrequest *BillRequest
}

func (suite *BillTestSuite) SetupTest() {
	suite.validBillrequest = &BillRequest{
		Description:  "bill 01",
		CustomerID:   "customer id",
		CurrencyCode: "USD",
		PeriodStart:  time.Now().UTC(),
		PeriodEnd:    time.Now().UTC(),
	}

	suite.invalidBillrequest = &BillRequest{
		Description: "bill 01",
	}
}

func (suite *BillTestSuite) Test_IsValidReturnFalse() {
	suite.False(suite.invalidBillrequest.IsValid())
}

func (suite *BillTestSuite) Test_IsValidReturnTrue() {
	suite.True(suite.validBillrequest.IsValid())
}

func TestBillTestSuite(t *testing.T) {
	suite.Run(t, new(BillTestSuite))
}
