package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserUUIDFromContext 從 Gin context 取得已認證的使用者 UUID
func UserUUIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	id, ok := userID.(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, false
	}

	return id, true
}
