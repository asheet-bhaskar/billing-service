package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CustomerTestSuite struct {
	suite.Suite
	validCustomer   *Customer
	invalidCustomer *Customer
}

func (suite *CustomerTestSuite) SetupTest() {
	suite.validCustomer = &Customer{
		FirstName: "John",
		LastName:  "jacobs",
		Email:     "john.jacobs@mail.com",
	}
	suite.invalidCustomer = &Customer{
		FirstName: "",
	}
}

func (suite *CustomerTestSuite) Test_IsValidReturnFalse() {
	suite.False(suite.invalidCustomer.IsValid())
}

func (suite *CustomerTestSuite) Test_IsValidReturnTrue() {
	suite.True(suite.validCustomer.IsValid())
}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}
