package test

import (
	"member_API/models"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testDB   *gorm.DB
	dbOnce   sync.Once
	dbPath   string
)

// SetupTestDB 設置共享的測試資料庫（整個測試套件共用一個）
func SetupTestDB(t *testing.T) *gorm.DB {
	dbOnce.Do(func() {
		tmpDir := os.TempDir()
		dbPath = filepath.Join(tmpDir, "shared_test.db")

		// 使用純 Go 的 SQLite 驅動（不需要 CGO）
		// 當 CGO_ENABLED=0 時，gorm.io/driver/sqlite 會自動使用純 Go 驅動（modernc.org/sqlite）
		// 因此不需要手動導入 modernc.org/sqlite
		var err error
		testDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v", err)
		}

		err = testDB.AutoMigrate(&models.Member{}, &models.Product{})
		if err != nil {
			t.Fatalf("Failed to migrate test database: %v", err)
		}
	})

	return testDB
}

// BeginTestTransaction 為單個測試開始一個 transaction
// 測試結束後會自動回滾，不需要手動清理資料
func BeginTestTransaction(t *testing.T, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}

// CleanupTestDB 清理測試資料庫（可在 TestMain 中呼叫）
func CleanupTestDB() {
	if testDB != nil {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		if dbPath != "" {
			_ = os.Remove(dbPath)
		}
	}
}

