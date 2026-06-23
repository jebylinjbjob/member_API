package auth

import (
	"net/http"
	"strings"

	"member_API/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 認證中間件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "缺少 Authorization header"})
			c.Abort()
			return
		}

		// 檢查 Bearer 前綴
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{response.KeyError: "Authorization header 格式錯誤，應為 'Bearer {token}'"},
			)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "無效的 token"})
			c.Abort()
			return
		}

		userID, err := claims.SubjectID()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{response.KeyError: "無效的 token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
