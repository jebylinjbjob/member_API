package services

import (
	"errors"
	"time"

	"member_API/auth"
	"member_API/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberService struct {
	DB *gorm.DB
}

func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{DB: db}
}

// CreateMember 建立新會員
func (s *MemberService) CreateMember(
	name, email, password string,
	creatorId uint,
) (*models.Member, error) {
	// 檢查 email 是否已存在
	var exists models.Member
	if err := s.DB.Where("email = ? AND is_deleted = ?", email, false).
		First(&exists).
		Error; err == nil {
		return nil, errors.New("email 已被使用")
	}

	// 加密密碼
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	member := &models.Member{
		Base: models.Base{
			CreationTime: now,
			CreatorId:    creatorId,
			IsDeleted:    false,
		},
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}

	if err := s.DB.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

// UpdateMember 更新會員資訊
func (s *MemberService) UpdateMember(
	id uint,
	name, email string,
	modifierId uint,
) (*models.Member, error) {
	var member models.Member
	if err := s.DB.Where("is_deleted = ?", false).First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("會員不存在")
		}
		return nil, err
	}

	now := time.Now()
	member.Name = name
	member.Email = email
	member.LastModificationTime = &now
	member.LastModifierId = modifierId

	if err := s.DB.Save(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// DeleteMember 軟刪除會員
func (s *MemberService) DeleteMember(id, deleterId uint) error {
	now := time.Now()
	result := s.DB.Model(&models.Member{}).
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

// GetMemberByEmail 依 email 取得單一會員
func (s *MemberService) GetMemberByEmail(email string) (*models.Member, error) {
	var member models.Member
	if err := s.DB.Where("email = ? AND is_deleted = ?", email, false).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("會員不存在")
		}
		return nil, err
	}
	return &member, nil
}

// ResetPassword 重設會員密碼
func (s *MemberService) ResetPassword(id uuid.UUID, newPassword string) error {
	member, err := s.GetMemberByUUID(id)
	if err != nil {
		return err
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.DB.Model(member).Updates(map[string]interface{}{
		"password_hash":          hash,
		"last_modification_time": &now,
	}).Error
}

// GetMemberByUUID 依 UUID 取得單一會員
func (s *MemberService) GetMemberByUUID(id uuid.UUID) (*models.Member, error) {
	var member models.Member
	if err := s.DB.Where("uuid = ? AND is_deleted = ?", id, false).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("會員不存在")
		}
		return nil, err
	}
	return &member, nil
}

// GetMemberByID 取得單一會員
func (s *MemberService) GetMemberByID(id uint) (*models.Member, error) {
	var member models.Member
	if err := s.DB.Where("is_deleted = ?", false).First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("會員不存在")
		}
		return nil, err
	}
	return &member, nil
}

// GetMembers 取得會員列表
func (s *MemberService) GetMembers(limit int) ([]models.Member, error) {
	var members []models.Member
	if err := s.DB.Where("is_deleted = ?", false).Limit(limit).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
