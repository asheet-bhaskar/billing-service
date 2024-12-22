package repository

import (
	"context"
	"testing"
	"time"

	models "github.com/asheet-bhaskar/billing-service/app/services/customer"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/stretchr/testify/suite"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	cr       customerRepository
	customer *models.Customer
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.cr = NewCustomerRepository(dbClient.DB)
	suite.customer = &models.Customer{
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacon@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (suite *CustomerRepositoryTestSuite) Test_CreateCustomerWhenSucceeds() {
	customer, err := suite.cr.Create(context.Background(), suite.customer)
	suite.Nil(err, "error should be nil")
	suite.NotNil(customer.ID)
}

func (suite *CustomerRepositoryTestSuite) Test_GetCustomerByIDWhenSucceeds() {
	customerRecord, err := suite.cr.GetByID(context.Background(), 1)

	suite.Nil(err, "error should be nil")
	suite.Equal(int64(1), customerRecord.ID)
	suite.Equal(suite.customer.FirstName, customerRecord.FirstName)
	suite.Equal(suite.customer.LastName, customerRecord.LastName)
	suite.Equal(suite.customer.Email, customerRecord.Email)
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
