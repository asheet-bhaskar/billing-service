package service

import (
	"context"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/stretchr/testify/mock"
)

type CustomerServiceMock struct {
	mock.Mock
}

func (m *CustomerServiceMock) Create(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *CustomerServiceMock) GetByID(ctx context.Context, id string) (*models.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Customer), args.Error(1)
}

type CurrencyServiceMock struct {
	mock.Mock
}

func (m *CurrencyServiceMock) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *CurrencyServiceMock) GetByID(ctx context.Context, id string) (*models.Currency, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Currency), args.Error(1)
}

type BillServiceMock struct {
	mock.Mock
}

func (m *BillServiceMock) Create(ctx context.Context, request *models.BillRequest) (*models.Bill, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *BillServiceMock) GetByID(ctx context.Context, id string) (*models.Bill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *BillServiceMock) AddLineItems(ctx context.Context, lineItem *models.LineItem) (*models.LineItem, error) {
	args := m.Called(ctx, lineItem)
	return args.Get(0).(*models.LineItem), args.Error(1)
}

func (m *BillServiceMock) RemoveLineItems(ctx context.Context, billID string, itemID string) (*models.LineItem, error) {
	args := m.Called(ctx, billID, itemID)
	return args.Get(0).(*models.LineItem), args.Error(1)
}

func (m *BillServiceMock) Close(ctx context.Context, id string) (*models.Bill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Bill), args.Error(1)
}

func (m *BillServiceMock) Invoice(ctx context.Context, billID string) (*models.Invoice, error) {
	args := m.Called(ctx, billID)
	return args.Get(0).(*models.Invoice), args.Error(1)
}
