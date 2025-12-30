package repositories

import (
	"member_API/models"

	"github.com/stretchr/testify/mock"
)

// MockMemberRepository 是 MemberRepository 的 Mock 實作
type MockMemberRepository struct {
	mock.Mock
}

func (m *MockMemberRepository) FindByEmail(email string) (*models.Member, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Member), args.Error(1)
}

func (m *MockMemberRepository) FindByID(id uint) (*models.Member, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Member), args.Error(1)
}

func (m *MockMemberRepository) Create(member *models.Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberRepository) Update(member *models.Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberRepository) SoftDelete(id uint, deleterId uint) error {
	args := m.Called(id, deleterId)
	return args.Error(0)
}

func (m *MockMemberRepository) FindAll(limit int) ([]models.Member, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Member), args.Error(1)
}

