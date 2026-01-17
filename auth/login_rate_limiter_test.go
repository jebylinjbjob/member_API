package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginRateLimiter_CheckAndRecord(t *testing.T) {
	limiter := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		MaxAttempts:   3,
		WindowPeriod:  1 * time.Minute,
		BlockDuration: 2 * time.Minute,
	}

	testIP := "192.168.1.1"

	// 第一次嘗試應該被允許
	allowed, _ := limiter.CheckAndRecord(testIP)
	assert.True(t, allowed, "First attempt should be allowed")

	// 第二次嘗試應該被允許
	allowed, _ = limiter.CheckAndRecord(testIP)
	assert.True(t, allowed, "Second attempt should be allowed")

	// 第三次嘗試應該被封鎖（因為 count 會變成 3，達到 MaxAttempts）
	allowed, remainingTime := limiter.CheckAndRecord(testIP)
	assert.False(t, allowed, "Third attempt should be blocked (count reaches MaxAttempts)")
	assert.Greater(t, remainingTime, time.Duration(0), "Should have remaining block time")

	// 再次嘗試應該仍然被封鎖
	allowed, _ = limiter.CheckAndRecord(testIP)
	assert.False(t, allowed, "Fourth attempt should still be blocked")
}

func TestLoginRateLimiter_Reset(t *testing.T) {
	limiter := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		MaxAttempts:   3,
		WindowPeriod:  1 * time.Minute,
		BlockDuration: 2 * time.Minute,
	}

	testIP := "192.168.1.2"

	// 進行兩次嘗試
	limiter.CheckAndRecord(testIP)
	limiter.CheckAndRecord(testIP)

	// 重置
	limiter.Reset(testIP)

	// 重置後應該可以重新開始計數
	allowed, _ := limiter.CheckAndRecord(testIP)
	assert.True(t, allowed, "After reset, should be allowed again")

	attempt, exists := limiter.attempts[testIP]
	assert.True(t, exists, "Should have new attempt record")
	assert.Equal(t, 1, attempt.Count, "Count should be 1 after reset")
}

func TestLoginRateLimiter_WindowExpiration(t *testing.T) {
	limiter := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		MaxAttempts:   3,
		WindowPeriod:  100 * time.Millisecond, // 短時間窗口便於測試
		BlockDuration: 2 * time.Minute,
	}

	testIP := "192.168.1.3"

	// 進行兩次嘗試
	limiter.CheckAndRecord(testIP)
	limiter.CheckAndRecord(testIP)

	// 等待時間窗口過期
	time.Sleep(150 * time.Millisecond)

	// 窗口過期後應該重置計數
	allowed, _ := limiter.CheckAndRecord(testIP)
	assert.True(t, allowed, "After window expiration, should be allowed")

	attempt := limiter.attempts[testIP]
	assert.Equal(t, 1, attempt.Count, "Count should reset to 1 after window expiration")
}

func TestLoginRateLimiter_BlockExpiration(t *testing.T) {
	limiter := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		MaxAttempts:   2,
		WindowPeriod:  1 * time.Minute,
		BlockDuration: 100 * time.Millisecond, // 短封鎖時間便於測試
	}

	testIP := "192.168.1.4"

	// 進行足夠次數的嘗試以觸發封鎖
	limiter.CheckAndRecord(testIP)
	limiter.CheckAndRecord(testIP)

	// 應該被封鎖
	allowed, _ := limiter.CheckAndRecord(testIP)
	assert.False(t, allowed, "Should be blocked after max attempts")

	// 等待封鎖過期
	time.Sleep(150 * time.Millisecond)

	// 封鎖過期後應該可以再次嘗試
	allowed, _ = limiter.CheckAndRecord(testIP)
	assert.True(t, allowed, "After block expiration, should be allowed")
}

func TestLoginRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建新的限制器用於測試
	limiter := &LoginRateLimiter{
		attempts:      make(map[string]*LoginAttempt),
		MaxAttempts:   3,
		WindowPeriod:  1 * time.Minute,
		BlockDuration: 2 * time.Minute,
	}

	router := gin.New()
	router.SetTrustedProxies(nil)
	
	// 創建一個使用測試限制器的中間件
	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed, remainingTime := limiter.CheckAndRecord(clientIP)
		
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "登入嘗試過於頻繁，請稍後再試",
				"retry_after_seconds": int(remainingTime.Seconds()),
			})
			c.Abort()
			return
		}
		
		c.Next()
	})
	
	router.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 進行多次請求，使用相同的遠端 IP
	testIP := "192.168.1.100"
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("POST", "/login", nil)
		req.RemoteAddr = testIP + ":12345"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// 第三次請求應該被封鎖（因為會達到 MaxAttempts=3）
	req := httptest.NewRequest("POST", "/login", nil)
	req.RemoteAddr = testIP + ":12345"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Third request should be blocked")
}

func TestGetLoginRateLimiter_Singleton(t *testing.T) {
	limiter1 := GetLoginRateLimiter()
	limiter2 := GetLoginRateLimiter()

	assert.Same(t, limiter1, limiter2, "Should return the same instance")
}
