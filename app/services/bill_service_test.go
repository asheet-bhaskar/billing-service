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

type MockBillRepository struct {
	mock.Mock
}

func (m *MockBillRepository) Create(ctx context.Context, bill *models.Bill) (*models.Bill, error) {
	args := m.Called(ctx, bill)
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *MockBillRepository) GetByID(ctx context.Context, id int64) (*models.Bill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Bill), args.Error(1)
}

type BillServiceTestSuite struct {
	suite.Suite
	MockRepo *MockBillRepository
	bs       BillService
	bill     *models.Bill
}

func (suite *BillServiceTestSuite) SetupTest() {

	mockRepo := new(MockBillRepository)
	suite.MockRepo = mockRepo
	suite.bs = NewBillService(mockRepo)

	suite.bill = &models.Bill{
		Description: "Bill - 01",
		CustomerID:  100,
		CurrencyID:  2,
		Status:      "open",
		TotalAmount: 100.45,
		PeriodStart: time.Now().UTC(),
		PeriodEnd:   time.Now().Add(time.Hour * 100),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsErrorWhenFails() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.bill).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.Create(ctx, suite.bill)

	suite.Require().Error(err)
	suite.Require().NotEqual(suite.bill, bill)
}

func (suite *BillServiceTestSuite) Test_CreateBillReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	suite.MockRepo.On("Create", ctx, suite.bill).Return(suite.bill, nil)

	bill, err := suite.bs.Create(ctx, suite.bill)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.bill, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsErrorWhenFails() {
	ctx := context.Background()
	billID := int64(1)
	suite.MockRepo.On("GetByID", ctx, billID).Return(&models.Bill{}, errors.New("test-error"))

	bill, err := suite.bs.GetByID(ctx, billID)

	suite.Require().NotNil(err)
	suite.Require().NotEqual(suite.bill, bill)
}

func (suite *BillServiceTestSuite) Test_GetByIDReturnsNilErrorWhenSucceeds() {
	ctx := context.Background()
	billID := int64(1)
	suite.MockRepo.On("GetByID", ctx, billID).Return(suite.bill, nil)

	bill, err := suite.bs.GetByID(ctx, billID)

	suite.Require().Nil(err)
	suite.Require().Equal(suite.bill, bill)
}

func TestBillServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
