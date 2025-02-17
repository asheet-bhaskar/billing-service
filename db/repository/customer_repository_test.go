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

type CustomerRepositoryTestSuite struct {
	suite.Suite
	dbClient *gorm.DB
	cr       CustomerRepository
	customer *models.Customer
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	host := "localhost"
	port := "5434"
	user := "billing_service_test"
	password := "billing_service_test"
	name := "billing_service_test"
	migrationsPath := "../migrations"

	dbClient, err := database.InitDBClient(host, port, user, password, name, migrationsPath)

	suite.Nil(err, "error should be nil")

	suite.dbClient = dbClient.DB

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

func (suite *CustomerRepositoryTestSuite) TearDownSuite() {
	fmt.Printf("cleaning up db records")
	suite.dbClient.Exec("DELETE FROM customers")
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
