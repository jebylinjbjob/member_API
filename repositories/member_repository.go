package repositories

import "member_API/models"

// MemberRepository 定義會員資料庫操作的接口
type MemberRepository interface {
	// FindByEmail 根據 email 查找會員（用於檢查重複）
	FindByEmail(email string) (*models.Member, error)

	// FindByID 根據 ID 查找會員
	FindByID(id uint) (*models.Member, error)

	// Create 創建會員
	Create(member *models.Member) error

	// Update 更新會員
	Update(member *models.Member) error

	// SoftDelete 軟刪除會員
	SoftDelete(id uint, deleterId uint) error

	// FindAll 取得會員列表
	FindAll(limit int) ([]models.Member, error)
}

