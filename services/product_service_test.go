package services

import (
	"errors"
	"member_API/models"
	"member_API/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_CreateProduct_Success(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.CreateProduct("iPhone 15", 35900.0, "最新款 iPhone", "image.jpg", 100, 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "iPhone 15", product.ProductName)
	assert.Equal(t, 35900.0, product.ProductPrice)
	assert.Equal(t, 100, product.ProductStock)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_Success(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProduct := &models.Product{
		Base:           models.Base{ID: 1},
		ProductName:    "iPhone 15",
		ProductPrice:   35900.0,
		ProductStock:   100,
	}
	mockRepo.On("FindByID", uint(1)).Return(expectedProduct, nil)

	product, err := service.GetProductByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "iPhone 15", product.ProductName)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	product, err := service.GetProductByID(999)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "產品不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_Success(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	existingProduct := &models.Product{
		Base:           models.Base{ID: 1},
		ProductName:    "iPhone 15",
		ProductPrice:   35900.0,
		ProductStock:   100,
	}
	updatedProduct := &models.Product{
		Base:           models.Base{ID: 1},
		ProductName:    "iPhone 15 Pro",
		ProductPrice:   42900.0,
		ProductStock:   50,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingProduct, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(updatedProduct, nil).Once()

	updates := map[string]interface{}{
		"product_name": "iPhone 15 Pro",
		"product_price": 42900.0,
		"product_stock": 50,
	}

	product, err := service.UpdateProduct(1, updates, 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "iPhone 15 Pro", product.ProductName)
	assert.Equal(t, 42900.0, product.ProductPrice)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	updates := map[string]interface{}{
		"product_name": "新名稱",
	}

	product, err := service.UpdateProduct(999, updates, 1)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "產品不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_Success(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("SoftDelete", uint(1), uint(1)).Return(nil)

	err := service.DeleteProduct(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("SoftDelete", uint(999), uint(1)).Return(errors.New("產品不存在或已被刪除"))

	err := service.DeleteProduct(999, 1)

	assert.Error(t, err)
	assert.Equal(t, "產品不存在或已被刪除", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProducts_Success(t *testing.T) {
	mockRepo := new(repositories.MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []models.Product{
		{Base: models.Base{ID: 1}, ProductName: "iPhone 15", ProductPrice: 35900.0},
		{Base: models.Base{ID: 2}, ProductName: "iPad Pro", ProductPrice: 27900.0},
	}
	mockRepo.On("Count").Return(int64(2), nil)
	mockRepo.On("FindAll", 50, 0).Return(expectedProducts, nil)

	products, total, err := service.GetProducts(50, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)
	assert.Equal(t, "iPhone 15", products[0].ProductName)
	mockRepo.AssertExpectations(t)
}

