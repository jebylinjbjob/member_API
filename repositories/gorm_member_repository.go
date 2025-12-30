package repositories

import (
	"errors"
	"member_API/models"
	"time"

	"gorm.io/gorm"
)

type gormMemberRepository struct {
	db *gorm.DB
}

// NewGormMemberRepository 創建 GORM MemberRepository 實例
func NewGormMemberRepository(db *gorm.DB) MemberRepository {
	return &gormMemberRepository{db: db}
}

func (r *gormMemberRepository) FindByEmail(email string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 沒找到，返回 nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *gormMemberRepository) FindByID(id uint) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("is_deleted = ?", false).First(&member, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *gormMemberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *gormMemberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *gormMemberRepository) SoftDelete(id uint, deleterId uint) error {
	now := time.Now()
	result := r.db.Model(&models.Member{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"is_deleted":             true,
			"deleted_at":             &now,
			"last_modifier_id":       deleterId,
			"last_modification_time": &now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("會員不存在或已被刪除")
	}
	return nil
}

func (r *gormMemberRepository) FindAll(limit int) ([]models.Member, error) {
	var members []models.Member
	err := r.db.Where("is_deleted = ?", false).Limit(limit).Find(&members).Error
	return members, err
}

