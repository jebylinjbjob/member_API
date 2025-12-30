//go:build integration

package services

import (
	"member_API/repositories"
	"member_API/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductService_CreateProduct_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	product, err := service.CreateProduct("iPhone 15", 35900.0, "最新款 iPhone", "image.jpg", 100, 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	if product != nil {
		assert.NotZero(t, product.ID)
		assert.Equal(t, "iPhone 15", product.ProductName)
		assert.Equal(t, 35900.0, product.ProductPrice)

		// 驗證可以從資料庫讀取
		retrieved, err := service.GetProductByID(product.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, product.ProductName, retrieved.ProductName)
	}
}

func TestProductService_GetProductByID_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 創建測試產品
	product, err := service.CreateProduct("測試產品", 1000.0, "測試描述", "test.jpg", 50, 1)
	assert.NoError(t, err)

	// 取得產品
	retrieved, err := service.GetProductByID(product.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, product.ID, retrieved.ID)
	assert.Equal(t, "測試產品", retrieved.ProductName)
}

func TestProductService_GetProductByID_NotFound_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 嘗試取得不存在的產品
	_, err := service.GetProductByID(99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "產品不存在")
}

func TestProductService_UpdateProduct_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 創建測試產品
	product, err := service.CreateProduct("原始產品", 1000.0, "原始描述", "old.jpg", 50, 1)
	assert.NoError(t, err)

	// 更新產品
	updates := map[string]interface{}{
		"product_name": "新產品",
		"product_price": 2000.0,
		"product_stock": 100,
	}
	updated, err := service.UpdateProduct(product.ID, updates, 1)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "新產品", updated.ProductName)
	assert.Equal(t, 2000.0, updated.ProductPrice)
	assert.Equal(t, 100, updated.ProductStock)

	// 驗證更新確實保存
	retrieved, err := service.GetProductByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "新產品", retrieved.ProductName)
	assert.Equal(t, 2000.0, retrieved.ProductPrice)
}

func TestProductService_DeleteProduct_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 創建測試產品
	product, err := service.CreateProduct("待刪除產品", 1000.0, "描述", "image.jpg", 50, 1)
	assert.NoError(t, err)

	// 刪除產品
	err = service.DeleteProduct(product.ID, 1)
	assert.NoError(t, err)

	// 驗證產品已被軟刪除（無法透過 GetProductByID 取得）
	_, err = service.GetProductByID(product.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "產品不存在")
}

func TestProductService_GetProducts_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 創建多個測試產品
	_, err := service.CreateProduct("產品1", 1000.0, "描述1", "image1.jpg", 50, 1)
	assert.NoError(t, err)
	_, err = service.CreateProduct("產品2", 2000.0, "描述2", "image2.jpg", 100, 1)
	assert.NoError(t, err)

	// 取得產品列表
	products, total, err := service.GetProducts(50, 0)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(total), 2)
	assert.GreaterOrEqual(t, len(products), 2)
}

func TestProductService_GetProducts_Pagination_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormProductRepository(tx)
	service := NewProductService(repo)

	// 創建多個測試產品
	for i := 0; i < 5; i++ {
		productName := "產品" + string(rune('A'+i))
		_, err := service.CreateProduct(
			productName,
			1000.0+float64(i)*100,
			"描述",
			"image.jpg",
			50,
			1,
		)
		assert.NoError(t, err)
	}

	// 測試分頁
	products1, total, err := service.GetProducts(2, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, products1, 2)

	products2, _, err := service.GetProducts(2, 2)
	assert.NoError(t, err)
	assert.Len(t, products2, 2)
	// 驗證分頁結果不同
	assert.NotEqual(t, products1[0].ID, products2[0].ID)
}

