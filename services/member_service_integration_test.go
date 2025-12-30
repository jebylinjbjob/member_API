//go:build integration

package services

import (
	"member_API/repositories"
	"member_API/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemberService_CreateMember_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	member, err := service.CreateMember("張三", "test@example.com", "password123", 1)

	// 驗證資料確實寫入資料庫
	assert.NoError(t, err)
	assert.NotNil(t, member)
	if member != nil {
		assert.NotZero(t, member.ID)
		assert.Equal(t, "test@example.com", member.Email)
		assert.Equal(t, "張三", member.Name)
		assert.NotEmpty(t, member.PasswordHash)

		// 驗證可以從資料庫讀取
		retrieved, err := service.GetMemberByID(member.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, member.Email, retrieved.Email)
	}
}

func TestMemberService_CreateMember_DuplicateEmail_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 創建第一個會員
	_, err := service.CreateMember("張三", "test@example.com", "password123", 1)
	assert.NoError(t, err)

	// 嘗試用相同 email 創建第二個會員
	_, err = service.CreateMember("李四", "test@example.com", "password456", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email 已被使用")
}

func TestMemberService_GetMemberByID_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 創建測試會員
	member, err := service.CreateMember("測試用戶", "test@example.com", "password123", 1)
	assert.NoError(t, err)

	// 取得會員
	retrieved, err := service.GetMemberByID(member.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, member.ID, retrieved.ID)
	assert.Equal(t, "測試用戶", retrieved.Name)
}

func TestMemberService_GetMemberByID_NotFound_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 嘗試取得不存在的會員
	_, err := service.GetMemberByID(99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "會員不存在")
}

func TestMemberService_UpdateMember_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 創建測試會員
	member, err := service.CreateMember("原始名稱", "original@example.com", "password123", 1)
	assert.NoError(t, err)

	// 更新會員
	updated, err := service.UpdateMember(member.ID, "新名稱", "new@example.com", 1)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "新名稱", updated.Name)
	assert.Equal(t, "new@example.com", updated.Email)

	// 驗證更新確實保存
	retrieved, err := service.GetMemberByID(member.ID)
	assert.NoError(t, err)
	assert.Equal(t, "新名稱", retrieved.Name)
	assert.Equal(t, "new@example.com", retrieved.Email)
}

func TestMemberService_DeleteMember_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 創建測試會員
	member, err := service.CreateMember("待刪除用戶", "delete@example.com", "password123", 1)
	assert.NoError(t, err)

	// 刪除會員
	err = service.DeleteMember(member.ID, 1)
	assert.NoError(t, err)

	// 驗證會員已被軟刪除（無法透過 GetMemberByID 取得）
	_, err = service.GetMemberByID(member.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "會員不存在")
}

func TestMemberService_GetMembers_Integration(t *testing.T) {
	db := test.SetupTestDB(t)
	tx := test.BeginTestTransaction(t, db)

	repo := repositories.NewGormMemberRepository(tx)
	service := NewMemberService(repo)

	// 創建多個測試會員
	_, err := service.CreateMember("用戶1", "user1@example.com", "password123", 1)
	assert.NoError(t, err)
	_, err = service.CreateMember("用戶2", "user2@example.com", "password123", 1)
	assert.NoError(t, err)

	// 取得會員列表
	members, err := service.GetMembers(50)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(members), 2)
}

