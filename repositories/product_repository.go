package repositories

import "member_API/models"

// ProductRepository 定義產品資料庫操作的接口
type ProductRepository interface {
	// FindByID 根據 ID 查找產品
	FindByID(id uint) (*models.Product, error)

	// Create 創建產品
	Create(product *models.Product) error

	// Update 更新產品
	Update(product *models.Product) error

	// SoftDelete 軟刪除產品
	SoftDelete(id uint, deleterId uint) error

	// FindAll 取得產品列表（支持分頁）
	FindAll(limit, offset int) ([]models.Product, error)

	// Count 取得產品總數（排除已刪除）
	Count() (int64, error)
}

