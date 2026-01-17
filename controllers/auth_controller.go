package controllers

import (
	"errors"
	"net/http"
	"time"

	"member_API/auth"
	"member_API/models"
	"member_API/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"張三"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// Register 用戶註冊
// @Summary 用戶註冊
// @Description 註冊新用戶，返回 JWT token 和用戶信息
// @Tags 認證
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "註冊信息"
// @Success 201 {object} AuthResponse "註冊成功"
// @Failure 400 {object} map[string]string "請求參數錯誤"
// @Failure 409 {object} map[string]string "該電子郵件已被註冊"
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /register [post]
func Register(input *gin.Context) {
	if db == nil {
		input.JSON(http.StatusInternalServerError, gin.H{"error": "數據庫連接未配置"})
		return
	}

	var req RegisterRequest
	if err := input.ShouldBindJSON(&req); err != nil {
		input.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 Service 層建立會員（自動處理密碼加密、審計欄位等）
	svc := services.NewMemberService(db)

	// 註冊時使用 creatorId = 0 表示自行註冊
	member, err := svc.CreateMember(req.Name, req.Email, req.Password, 0)
	if err != nil {
		if err.Error() == "email 已被使用" {
			input.JSON(http.StatusConflict, gin.H{"error": "該電子郵件已被註冊"})
			return
		}
		input.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := User{ID: int64(member.ID), Name: member.Name, Email: member.Email}

	// 生成 token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		input.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失敗"})
		return
	}

	input.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login 用戶登入
// @Summary 用戶登入
// @Description 用戶登入，驗證郵件和密碼後返回 JWT token 和用戶信息
// @Tags 認證
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登入信息"
// @Success 200 {object} AuthResponse "登入成功"
// @Failure 400 {object} map[string]string "請求參數錯誤"
// @Failure 401 {object} map[string]string "電子郵件或密碼錯誤"
// @Failure 403 {object} map[string]string "帳號已被鎖定"
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /login [post]
func Login(input *gin.Context) {
	if db == nil {
		input.JSON(http.StatusInternalServerError, gin.H{"error": "數據庫連接未配置"})
		return
	}

	var req LoginRequest
	if err := input.ShouldBindJSON(&req); err != nil {
		input.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查詢用戶
	var member models.Member
	err := db.WithContext(input.Request.Context()).
		Where("email = ?", req.Email).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			input.JSON(http.StatusUnauthorized, gin.H{"error": "電子郵件或密碼錯誤"})
			return
		}
		input.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 檢查帳號是否被鎖定
	if member.IsLocked && member.LockedUntil != nil && member.LockedUntil.After(time.Now()) {
		input.JSON(http.StatusForbidden, gin.H{"error": "帳號已被鎖定，請稍後再試"})
		return
	}

	// 如果鎖定時間已過，解鎖帳號
	if member.IsLocked && member.LockedUntil != nil && member.LockedUntil.Before(time.Now()) {
		db.Model(&member).Updates(map[string]interface{}{
			"is_locked":             false,
			"failed_login_attempts": 0,
			"locked_until":          nil,
		})
		member.IsLocked = false
		member.FailedLoginAttempts = 0
		member.LockedUntil = nil
	}

	// 驗證密碼
	if !auth.CheckPassword(req.Password, member.PasswordHash) {
		// 增加失敗登入次數
		member.FailedLoginAttempts++
		
		// 如果失敗次數達到5次，鎖定帳號30分鐘
		if member.FailedLoginAttempts >= 5 {
			lockUntil := time.Now().Add(30 * time.Minute)
			db.Model(&member).Updates(map[string]interface{}{
				"is_locked":             true,
				"failed_login_attempts": member.FailedLoginAttempts,
				"locked_until":          lockUntil,
			})
			input.JSON(http.StatusForbidden, gin.H{"error": "登入失敗次數過多，帳號已被鎖定30分鐘"})
			return
		}
		
		// 更新失敗次數
		db.Model(&member).Update("failed_login_attempts", member.FailedLoginAttempts)
		input.JSON(http.StatusUnauthorized, gin.H{"error": "電子郵件或密碼錯誤"})
		return
	}

	// 登入成功，重置失敗次數
	if member.FailedLoginAttempts > 0 {
		db.Model(&member).Update("failed_login_attempts", 0)
	}

	user := User{ID: int64(member.ID), Name: member.Name, Email: member.Email}

	// 生成 token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		input.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失敗"})
		return
	}

	input.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile 獲取當前用戶信息（需要認證）
// @Summary 獲取當前用戶信息
// @Description 獲取當前登入用戶的詳細信息，需要 JWT 認證
// @Tags 用戶
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]User "獲取成功"
// @Failure 401 {object} map[string]string "未認證"
// @Failure 404 {object} map[string]string "用戶不存在"
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未認證"})
		return
	}

	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "數據庫連接未配置"})
		return
	}

	idValue, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未認證"})
		return
	}

	var member models.Member
	if err := db.WithContext(c.Request.Context()).
		Select("id", "name", "email").
		First(&member, idValue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用戶不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": User{ID: int64(member.ID), Name: member.Name, Email: member.Email}})
}
