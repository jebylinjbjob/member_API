package services

import (
	"errors"
	"member_API/models"
	"member_API/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemberService_CreateMember_Success(t *testing.T) {
	// 創建 Mock Repository
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	// 設置 Mock 預期行為
	// 1. FindByEmail 返回 nil（表示 email 不存在）
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

	// 2. Create 成功
	mockRepo.On("Create", mock.AnythingOfType("*models.Member")).Return(nil)

	// 執行測試
	member, err := service.CreateMember("張三", "test@example.com", "password123", 1)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "test@example.com", member.Email)
	assert.Equal(t, "張三", member.Name)
	assert.NotEmpty(t, member.PasswordHash) // 密碼應該被加密

	// 驗證所有 Mock 都被調用
	mockRepo.AssertExpectations(t)
}

func TestMemberService_CreateMember_EmailExists(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	// Mock：email 已存在
	existingMember := &models.Member{
		Base: models.Base{ID: 1},
		Email: "test@example.com",
	}
	mockRepo.On("FindByEmail", "test@example.com").Return(existingMember, nil)

	// 執行測試
	member, err := service.CreateMember("李四", "test@example.com", "password123", 1)

	// 驗證錯誤
	assert.Error(t, err)
	assert.Nil(t, member)
	assert.Equal(t, "email 已被使用", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMemberService_GetMemberByID_Success(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	expectedMember := &models.Member{
		Base: models.Base{ID: 1},
		Name:  "張三",
		Email: "test@example.com",
	}
	mockRepo.On("FindByID", uint(1)).Return(expectedMember, nil)

	member, err := service.GetMemberByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, uint(1), member.ID)
	assert.Equal(t, "張三", member.Name)
	mockRepo.AssertExpectations(t)
}

func TestMemberService_GetMemberByID_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	member, err := service.GetMemberByID(999)

	assert.Error(t, err)
	assert.Nil(t, member)
	assert.Equal(t, "會員不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMemberService_UpdateMember_Success(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	existingMember := &models.Member{
		Base: models.Base{ID: 1},
		Name:  "原始名稱",
		Email: "original@example.com",
	}
	mockRepo.On("FindByID", uint(1)).Return(existingMember, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Member")).Return(nil)

	updated, err := service.UpdateMember(1, "新名稱", "new@example.com", 1)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "新名稱", updated.Name)
	assert.Equal(t, "new@example.com", updated.Email)
	assert.NotNil(t, updated.LastModificationTime)
	mockRepo.AssertExpectations(t)
}

func TestMemberService_UpdateMember_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	updated, err := service.UpdateMember(999, "新名稱", "new@example.com", 1)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, "會員不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMemberService_DeleteMember_Success(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	mockRepo.On("SoftDelete", uint(1), uint(1)).Return(nil)

	err := service.DeleteMember(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMemberService_DeleteMember_NotFound(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	mockRepo.On("SoftDelete", uint(999), uint(1)).Return(errors.New("會員不存在或已被刪除"))

	err := service.DeleteMember(999, 1)

	assert.Error(t, err)
	assert.Equal(t, "會員不存在或已被刪除", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMemberService_GetMembers_Success(t *testing.T) {
	mockRepo := new(repositories.MockMemberRepository)
	service := NewMemberService(mockRepo)

	expectedMembers := []models.Member{
		{Base: models.Base{ID: 1}, Name: "張三", Email: "test1@example.com"},
		{Base: models.Base{ID: 2}, Name: "李四", Email: "test2@example.com"},
	}
	mockRepo.On("FindAll", 50).Return(expectedMembers, nil)

	members, err := service.GetMembers(50)

	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, "張三", members[0].Name)
	mockRepo.AssertExpectations(t)
}

