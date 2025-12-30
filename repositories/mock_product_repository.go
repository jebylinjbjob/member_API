package repositories

import (
	"member_API/models"

	"github.com/stretchr/testify/mock"
)

// MockProductRepository 是 ProductRepository 的 Mock 實作
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) SoftDelete(id uint, deleterId uint) error {
	args := m.Called(id, deleterId)
	return args.Error(0)
}

func (m *MockProductRepository) FindAll(limit, offset int) ([]models.Product, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

