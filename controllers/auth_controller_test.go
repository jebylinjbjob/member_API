package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"member_API/auth"
	"member_API/models"
	"member_API/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin_AccountLocking(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer func() {
		if err := testutil.CleanupTestDB(db); err != nil {
			t.Errorf("Failed to cleanup test database: %v", err)
		}
	}()

	SetupUserController(db)

	// Create a test user
	passwordHash, _ := auth.HashPassword("password123")
	member := &models.Member{
		Base: models.Base{
			CreationTime: time.Now(),
			CreatorId:    0,
			IsDeleted:    false,
		},
		Name:                "Test User",
		Email:               "test@example.com",
		PasswordHash:        passwordHash,
		IsLocked:            false,
		FailedLoginAttempts: 0,
	}
	db.Create(member)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", Login)

	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectedError  string
		setupFunc      func()
	}{
		{
			name:           "Successful login",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
			expectedError:  "",
			setupFunc: func() {
				db.Model(member).Updates(map[string]interface{}{
					"is_locked":             false,
					"failed_login_attempts": 0,
					"locked_until":          nil,
				})
			},
		},
		{
			name:           "Failed login - wrong password (1st attempt)",
			email:          "test@example.com",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "電子郵件或密碼錯誤",
			setupFunc: func() {
				db.Model(member).Updates(map[string]interface{}{
					"is_locked":             false,
					"failed_login_attempts": 0,
					"locked_until":          nil,
				})
			},
		},
		{
			name:           "Account locked after 5 failed attempts",
			email:          "test@example.com",
			password:       "wrongpassword",
			expectedStatus: http.StatusForbidden,
			expectedError:  "登入失敗次數過多，帳號已被鎖定30分鐘",
			setupFunc: func() {
				db.Model(member).Updates(map[string]interface{}{
					"is_locked":             false,
					"failed_login_attempts": 4,
					"locked_until":          nil,
				})
			},
		},
		{
			name:           "Login blocked when account is locked",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusForbidden,
			expectedError:  "帳號已被鎖定，請稍後再試",
			setupFunc: func() {
				lockUntil := time.Now().Add(30 * time.Minute)
				db.Model(member).Updates(map[string]interface{}{
					"is_locked":             true,
					"failed_login_attempts": 5,
					"locked_until":          lockUntil,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			loginReq := LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}
			body, _ := json.Marshal(loginReq)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestLogin_AutoUnlock(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer func() {
		if err := testutil.CleanupTestDB(db); err != nil {
			t.Errorf("Failed to cleanup test database: %v", err)
		}
	}()

	SetupUserController(db)

	// Create a test user with expired lock
	passwordHash, _ := auth.HashPassword("password123")
	lockUntil := time.Now().Add(-1 * time.Minute) // Lock expired 1 minute ago
	member := &models.Member{
		Base: models.Base{
			CreationTime: time.Now(),
			CreatorId:    0,
			IsDeleted:    false,
		},
		Name:                "Test User",
		Email:               "test2@example.com",
		PasswordHash:        passwordHash,
		IsLocked:            true,
		FailedLoginAttempts: 5,
		LockedUntil:         &lockUntil,
	}
	db.Create(member)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", Login)

	loginReq := LoginRequest{
		Email:    "test2@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed because lock has expired
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify member is unlocked in database
	var updatedMember models.Member
	db.First(&updatedMember, member.ID)
	assert.False(t, updatedMember.IsLocked)
	assert.Equal(t, 0, updatedMember.FailedLoginAttempts)
}
