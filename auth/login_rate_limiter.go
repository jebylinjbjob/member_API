package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginAttempt 記錄登入嘗試
type LoginAttempt struct {
	Count      int
	FirstAttempt time.Time
	BlockedUntil *time.Time
}

// LoginRateLimiter 登入速率限制器
type LoginRateLimiter struct {
	attempts map[string]*LoginAttempt
	mu       sync.RWMutex
	
	// 配置參數
	MaxAttempts   int           // 最大嘗試次數
	WindowPeriod  time.Duration // 時間窗口
	BlockDuration time.Duration // 封鎖時長
}

var (
	// 全局登入速率限制器
	globalLoginLimiter *LoginRateLimiter
	limiterOnce        sync.Once
)

// GetLoginRateLimiter 獲取全局登入速率限制器（單例模式）
func GetLoginRateLimiter() *LoginRateLimiter {
	limiterOnce.Do(func() {
		globalLoginLimiter = &LoginRateLimiter{
			attempts:      make(map[string]*LoginAttempt),
			MaxAttempts:   10,                // 10次嘗試
			WindowPeriod:  5 * time.Minute,   // 5分鐘內
			BlockDuration: 15 * time.Minute,  // 封鎖15分鐘
		}
		// 啟動清理協程
		go globalLoginLimiter.cleanup()
	})
	return globalLoginLimiter
}

// cleanup 定期清理過期記錄
func (l *LoginRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, attempt := range l.attempts {
			// 清理已解鎖且超過窗口期的記錄
			if attempt.BlockedUntil == nil && now.Sub(attempt.FirstAttempt) > l.WindowPeriod {
				delete(l.attempts, key)
			}
			// 清理已過期的封鎖記錄
			if attempt.BlockedUntil != nil && now.After(*attempt.BlockedUntil) {
				delete(l.attempts, key)
			}
		}
		l.mu.Unlock()
	}
}

// CheckAndRecord 檢查並記錄登入嘗試
// 返回 true 表示允許嘗試，false 表示被封鎖
func (l *LoginRateLimiter) CheckAndRecord(key string) (allowed bool, remainingTime time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	now := time.Now()
	attempt, exists := l.attempts[key]
	
	if !exists {
		// 首次嘗試
		l.attempts[key] = &LoginAttempt{
			Count:        1,
			FirstAttempt: now,
		}
		return true, 0
	}
	
	// 檢查是否被封鎖
	if attempt.BlockedUntil != nil {
		if now.Before(*attempt.BlockedUntil) {
			return false, attempt.BlockedUntil.Sub(now)
		}
		// 封鎖已過期，重置
		delete(l.attempts, key)
		l.attempts[key] = &LoginAttempt{
			Count:        1,
			FirstAttempt: now,
		}
		return true, 0
	}
	
	// 檢查是否在時間窗口內
	if now.Sub(attempt.FirstAttempt) > l.WindowPeriod {
		// 超過時間窗口，重置計數
		l.attempts[key] = &LoginAttempt{
			Count:        1,
			FirstAttempt: now,
		}
		return true, 0
	}
	
	// 在時間窗口內，增加計數
	attempt.Count++
	
	// 檢查是否超過最大嘗試次數
	if attempt.Count >= l.MaxAttempts {
		blockedUntil := now.Add(l.BlockDuration)
		attempt.BlockedUntil = &blockedUntil
		return false, l.BlockDuration
	}
	
	return true, 0
}

// Reset 重置指定 key 的記錄（登入成功時調用）
func (l *LoginRateLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

// LoginRateLimitMiddleware 登入速率限制中間件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	limiter := GetLoginRateLimiter()
	
	return func(c *gin.Context) {
		// 獲取客戶端 IP
		clientIP := c.ClientIP()
		
		// 檢查是否允許嘗試
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
	}
}
