package services

import (
	"errors"
	"member_API/auth"
	"member_API/models"
	"member_API/repositories"
	"time"
)

type MemberService struct {
	repo repositories.MemberRepository
}

func NewMemberService(repo repositories.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

// CreateMember 建立新會員
func (s *MemberService) CreateMember(name, email, password string, creatorId uint) (*models.Member, error) {
	// 檢查 email 是否已存在
	exists, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists != nil {
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

	if err := s.repo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

// UpdateMember 更新會員資訊
func (s *MemberService) UpdateMember(id uint, name, email string, modifierId uint) (*models.Member, error) {
	member, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("會員不存在")
	}

	now := time.Now()
	member.Name = name
	member.Email = email
	member.LastModificationTime = &now
	member.LastModifierId = modifierId

	if err := s.repo.Update(member); err != nil {
		return nil, err
	}

	return member, nil
}

// DeleteMember 軟刪除會員
func (s *MemberService) DeleteMember(id uint, deleterId uint) error {
	return s.repo.SoftDelete(id, deleterId)
}

// GetMemberByID 取得單一會員
func (s *MemberService) GetMemberByID(id uint) (*models.Member, error) {
	member, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("會員不存在")
	}
	return member, nil
}

// GetMembers 取得會員列表
func (s *MemberService) GetMembers(limit int) ([]models.Member, error) {
	return s.repo.FindAll(limit)
}
