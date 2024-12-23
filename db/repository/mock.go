package repository

import (
	"context"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/stretchr/testify/mock"
)

type MockBillRepository struct {
	mock.Mock
}

type MockCurrencyRepository struct {
	mock.Mock
}

type MockCustomerRepository struct {
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

func (m *MockBillRepository) AddLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	args := m.Called(ctx, lineItem)
	return args.Get(0).(*models.LineItem), args.Error(1)
}
func (m *MockBillRepository) RemoveLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	args := m.Called(ctx, lineItem)
	return args.Get(0).(*models.LineItem), args.Error(1)
}
func (m *MockBillRepository) GetLineItemsByBillID(ctx context.Context, id int64) ([]*models.LineItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*models.LineItem), args.Error(1)
}
func (m *MockBillRepository) Close(ctx context.Context, id int64) (*models.Bill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *MockCurrencyRepository) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) GetByID(ctx context.Context, id int64) (*models.Currency, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) GetByCode(ctx context.Context, code string) (*models.Currency, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Customer), args.Error(1)
}
