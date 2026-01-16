# 測試文件 (Testing Documentation)

## 概述 (Overview)

本專案使用 Go 原生的 `testing` 套件進行單元測試，確保代碼的可靠性和可維護性。

## 測試架構 (Test Architecture)

### 已實現的測試 (Implemented Tests)

#### 1. Auth 套件測試 (`auth/`)
- **password_test.go** - 密碼加密與驗證測試
  - ✅ `TestHashPassword` - 測試密碼加密功能
  - ✅ `TestCheckPassword` - 測試密碼驗證功能
  - ✅ `TestHashPasswordConsistency` - 測試加密一致性（鹽值隨機性）

- **jwt_test.go** - JWT Token 生成與驗證測試
  - ✅ `TestGenerateToken` - 測試 Token 生成
  - ✅ `TestValidateToken` - 測試 Token 驗證
  - ✅ `TestValidateTokenExpired` - 測試過期 Token 驗證
  - ✅ `TestValidateTokenWrongSecret` - 測試錯誤密鑰驗證
  - ✅ `TestTokenRoundTrip` - 測試 Token 完整流程

#### 2. Services 套件測試 (`services/`)
- **member_service_test.go** - 會員服務測試
  - ✅ `TestMemberService_CreateMember` - 建立會員
  - ✅ `TestMemberService_GetMemberByID` - 取得單一會員
  - ✅ `TestMemberService_GetMembers` - 取得會員列表
  - ✅ `TestMemberService_UpdateMember` - 更新會員資訊
  - ✅ `TestMemberService_DeleteMember` - 刪除會員（軟刪除）

- **product_service_test.go** - 產品服務測試
  - ✅ `TestProductService_CreateProduct` - 建立產品
  - ✅ `TestProductService_GetProductByID` - 取得單一產品
  - ✅ `TestProductService_GetProducts` - 取得產品列表（含分頁）
  - ✅ `TestProductService_UpdateProduct` - 更新產品資訊
  - ✅ `TestProductService_DeleteProduct` - 刪除產品（軟刪除）

#### 3. 測試輔助工具 (`testutil/`)
- **testutil.go** - 測試輔助函數
  - `SetupTestDB()` - 建立記憶體內 SQLite 測試資料庫
  - `CleanupTestDB()` - 清理測試資料庫
  - `CreateTestMember()` - 建立測試會員
  - `CreateTestProduct()` - 建立測試產品

## 執行測試 (Running Tests)

### 執行所有測試
```bash
go test ./...
```

### 執行特定套件測試
```bash
# 測試 auth 套件
go test ./auth/...

# 測試 services 套件
go test ./services/...
```

### 執行測試並顯示詳細輸出
```bash
go test ./... -v
```

### 執行測試並顯示覆蓋率
```bash
go test ./... -cover
```

### 生成覆蓋率報告
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 測試策略 (Testing Strategy)

### 單元測試原則
1. **獨立性** - 每個測試應該獨立運行，不依賴其他測試
2. **可重複性** - 測試結果應該可重複且一致
3. **快速執行** - 使用記憶體資料庫以確保快速執行
4. **清晰命名** - 測試名稱清楚描述測試內容
5. **覆蓋率** - 目標達到 80% 以上的代碼覆蓋率

### 測試資料管理
- 使用 SQLite 記憶體資料庫 (`:memory:`) 進行測試
- 每個測試使用獨立的資料庫實例
- 測試結束後自動清理資料庫連接

### 測試案例設計
- **正常情境** - 測試正確的輸入和預期行為
- **邊界情境** - 測試邊界值和特殊情況
- **錯誤情境** - 測試錯誤處理和異常情況

## 未來測試計畫 (Future Testing Plan)

### Phase 1: 擴展單元測試 (已完成 ✅)
- [x] Auth 套件測試
- [x] Services 套件測試
- [x] 測試輔助工具建立

### Phase 2: Controller 層測試 (待實現)
- [ ] **auth_controller_test.go**
  - [ ] 測試用戶註冊端點
  - [ ] 測試用戶登入端點
  - [ ] 測試錯誤處理

- [ ] **user_controller_test.go**
  - [ ] 測試獲取所有用戶
  - [ ] 測試獲取單一用戶
  - [ ] 測試獲取當前用戶資訊
  - [ ] 測試 JWT 認證中介層

- [ ] **product_controller_test.go**
  - [ ] 測試產品 CRUD 操作
  - [ ] 測試分頁功能
  - [ ] 測試權限控制

### Phase 3: 整合測試 (Integration Tests)
- [ ] **API 端點整合測試**
  - [ ] 測試完整的認證流程（註冊 → 登入 → 獲取資訊）
  - [ ] 測試會員管理完整流程
  - [ ] 測試產品管理完整流程
  - [ ] 測試錯誤處理和邊界情況

- [ ] **GraphQL 整合測試**
  - [ ] 測試 GraphQL Queries
  - [ ] 測試 GraphQL Mutations
  - [ ] 測試 GraphQL 認證

- [ ] **資料庫整合測試**
  - [ ] 測試與真實 PostgreSQL 的互動
  - [ ] 測試資料庫遷移
  - [ ] 測試交易處理

### Phase 4: 端對端測試 (E2E Tests)
- [ ] 使用真實的 HTTP 請求測試完整流程
- [ ] 測試多用戶併發場景
- [ ] 測試錯誤恢復機制

### Phase 5: 效能測試 (Performance Tests)
- [ ] **負載測試**
  - [ ] API 端點負載測試
  - [ ] 資料庫查詢效能測試
  - [ ] 併發請求測試

- [ ] **壓力測試**
  - [ ] 系統極限測試
  - [ ] 資源使用監控

- [ ] **基準測試 (Benchmarks)**
  - [ ] 關鍵函數效能基準
  - [ ] 資料庫操作基準
  - [ ] JWT 生成與驗證基準

### Phase 6: 安全測試 (Security Tests)
- [ ] SQL 注入防護測試
- [ ] XSS 攻擊防護測試
- [ ] CSRF 防護測試
- [ ] JWT Token 安全性測試
- [ ] 密碼強度測試
- [ ] 權限控制測試

## 測試覆蓋率目標 (Coverage Goals)

| 套件 | 目標覆蓋率 | 當前狀態 |
|------|-----------|---------|
| auth | 80%+ | ✅ 已實現 |
| services | 80%+ | ✅ 已實現 |
| controllers | 80%+ | 🔄 待實現 |
| models | 60%+ | 📋 計畫中 |
| graphql | 70%+ | 📋 計畫中 |

## 持續整合 (Continuous Integration)

### 建議的 CI/CD 工作流程
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.24'
      - name: Run tests
        run: go test ./... -cover
      - name: Check coverage
        run: |
          go test ./... -coverprofile=coverage.out
          go tool cover -func=coverage.out
```

## 測試最佳實踐 (Testing Best Practices)

### 1. 測試命名規範
```go
func TestServiceName_MethodName(t *testing.T) {
    // 測試邏輯
}
```

### 2. 使用表格驅動測試
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"case1", "input1", "output1", false},
    {"case2", "input2", "output2", false},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 測試邏輯
    })
}
```

### 3. 使用子測試
```go
t.Run("描述性測試名稱", func(t *testing.T) {
    // 測試邏輯
})
```

### 4. 清理資源
```go
db := testutil.SetupTestDB(t)
defer testutil.CleanupTestDB(db)
```

## 除錯測試 (Debugging Tests)

### 執行單一測試
```bash
go test -run TestName ./path/to/package
```

### 顯示測試輸出
```bash
go test -v ./...
```

### 執行失敗時顯示完整堆疊
```bash
go test -v -failfast ./...
```

## 相關資源 (Resources)

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Advanced Testing with Go](https://about.sourcegraph.com/go/advanced-testing-in-go)
- [GORM Testing](https://gorm.io/docs/testing.html)

## 貢獻指南 (Contributing)

在提交 Pull Request 之前：
1. ✅ 確保所有測試通過：`go test ./...`
2. ✅ 確保代碼格式正確：`go fmt ./...`
3. ✅ 執行代碼檢查：`go vet ./...`
4. ✅ 為新功能添加相應的測試
5. ✅ 更新測試覆蓋率報告

## 問題回報 (Issue Reporting)

如果發現測試問題，請提供：
- 測試失敗的完整輸出
- 運行環境資訊（Go 版本、OS 等）
- 重現步驟
- 預期行為 vs 實際行為
