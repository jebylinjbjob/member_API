package controllers

import (
	"errors"
	"log"
	"net/http"

	"member_API/auth"
	"member_API/models"
	"member_API/response"
	"member_API/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"       example:"張三"`
	Email    string `json:"email"    binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"    binding:"required"       example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Password string `json:"password" binding:"required,min=6" example:"newpassword123"`
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
		input.JSON(
			http.StatusInternalServerError,
			gin.H{response.KeyError: response.MsgDBNotConfiguredZH},
		)
		return
	}

	var req RegisterRequest
	if err := input.ShouldBindJSON(&req); err != nil {
		input.JSON(http.StatusBadRequest, gin.H{response.KeyError: err.Error()})
		return
	}

	// 使用 Service 層建立會員（自動處理密碼加密、審計欄位等）
	svc := services.NewMemberService(db)

	// 註冊時使用 creatorId = 0 表示自行註冊
	member, err := svc.CreateMember(req.Name, req.Email, req.Password, 0)
	if err != nil {
		if err.Error() == "email 已被使用" {
			input.JSON(http.StatusConflict, gin.H{response.KeyError: "該電子郵件已被註冊"})
			return
		}
		input.JSON(http.StatusInternalServerError, gin.H{response.KeyError: err.Error()})
		return
	}

	user := User{ID: member.UUID.String(), Name: member.Name, Email: member.Email}

	token, err := auth.GenerateToken(member.UUID, "")
	if err != nil {
		input.JSON(http.StatusInternalServerError, gin.H{response.KeyError: "Token 生成失敗"})
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
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /login [post]
func Login(input *gin.Context) {
	if db == nil {
		input.JSON(
			http.StatusInternalServerError,
			gin.H{response.KeyError: response.MsgDBNotConfiguredZH},
		)
		return
	}

	var req LoginRequest
	if err := input.ShouldBindJSON(&req); err != nil {
		input.JSON(http.StatusBadRequest, gin.H{response.KeyError: err.Error()})
		return
	}

	// 查詢用戶
	var member models.Member
	err := db.WithContext(input.Request.Context()).
		Where("email = ?", req.Email).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			input.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "電子郵件或密碼錯誤"})
			return
		}
		input.JSON(http.StatusInternalServerError, gin.H{response.KeyError: err.Error()})
		return
	}

	// 驗證密碼
	if !auth.CheckPassword(req.Password, member.PasswordHash) {
		input.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "電子郵件或密碼錯誤"})
		return
	}

	user := User{ID: member.UUID.String(), Name: member.Name, Email: member.Email}

	token, err := auth.GenerateToken(member.UUID, "")
	if err != nil {
		input.JSON(http.StatusInternalServerError, gin.H{response.KeyError: "Token 生成失敗"})
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
	userUUID, ok := auth.UserUUIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "未認證"})
		return
	}

	if db == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{response.KeyError: response.MsgDBNotConfiguredZH},
		)
		return
	}

	var member models.Member
	if err := db.WithContext(c.Request.Context()).
		Select("id", "uuid", "name", "email").
		Where("uuid = ?", userUUID).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{response.KeyError: "用戶不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{response.KeyError: err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"user": User{ID: member.UUID.String(), Name: member.Name, Email: member.Email}},
	)
}

// ForgotPassword 請求重設密碼（對應 fastapi-users POST /auth/forgot-password）
// @Summary 忘記密碼
// @Description 請求重設密碼，無論 email 是否存在皆回傳 202 Accepted
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "電子郵件"
// @Success 202 "已接受請求"
// @Failure 400 {object} map[string]string "請求參數錯誤"
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	if db == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{response.KeyError: response.MsgDBNotConfiguredZH},
		)
		return
	}

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{response.KeyError: err.Error()})
		return
	}

	svc := services.NewMemberService(db)
	member, err := svc.GetMemberByEmail(req.Email)
	if err == nil {
		token, tokenErr := auth.GenerateResetPasswordToken(member.UUID)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{response.KeyError: "Token 生成失敗"})
			return
		}
		// 對應 hopenet on_after_forgot_password：開發階段記錄 token，正式環境可改為寄信
		log.Printf("User %s has forgot their password. Reset token: %s", member.UUID, token)
	}

	c.Status(http.StatusAccepted)
}

// ResetPassword 使用 token 重設密碼（對應 fastapi-users POST /auth/reset-password）
// @Summary 重設密碼
// @Description 使用 forgot-password 取得的 token 重設密碼
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "重設密碼資訊"
// @Success 200 "重設成功"
// @Failure 400 {object} map[string]interface{} "token 無效或密碼不符合規則"
// @Failure 500 {object} map[string]string "服務器錯誤"
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	if db == nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{response.KeyError: response.MsgDBNotConfiguredZH},
		)
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": gin.H{
				"code":   "RESET_PASSWORD_INVALID_PASSWORD",
				"reason": err.Error(),
			},
		})
		return
	}

	userUUID, err := auth.ValidateResetPasswordToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": auth.ErrResetPasswordBadToken.Error()})
		return
	}

	svc := services.NewMemberService(db)
	if err := svc.ResetPassword(userUUID, req.Password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "會員不存在" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": auth.ErrResetPasswordBadToken.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{response.KeyError: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
