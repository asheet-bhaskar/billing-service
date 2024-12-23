package repository

import (
	"context"
	"testing"
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	database "github.com/asheet-bhaskar/billing-service/db"
	"github.com/asheet-bhaskar/billing-service/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	cr       CustomerRepository
	customer *models.Customer
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	dbClient, err := database.InitDBClient()
	suite.Nil(err, "error should be nil")

	suite.cr = NewCustomerRepository(dbClient.DB)
	suite.customer = &models.Customer{
		ID:        utils.GetNewUUID(),
		FirstName: "John",
		LastName:  "Jacobs",
		Email:     "john.jacon@mail.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (suite *CustomerRepositoryTestSuite) Test_CreateCustomerWhenSucceeds() {
	_, err := suite.cr.Create(context.Background(), suite.customer)
	suite.Nil(err, "error should be nil")
}

func (suite *CustomerRepositoryTestSuite) Test_GetCustomerByIDWhenSucceeds() {
	_, err := suite.cr.Create(context.Background(), suite.customer)
	suite.Nil(err, "error should be nil")

	customerRecord, err := suite.cr.GetByID(context.Background(), suite.customer.ID)

	suite.Nil(err, "error should be nil")
	suite.Equal(suite.customer.FirstName, customerRecord.FirstName)
	suite.Equal(suite.customer.LastName, customerRecord.LastName)
	suite.Equal(suite.customer.Email, customerRecord.Email)
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
