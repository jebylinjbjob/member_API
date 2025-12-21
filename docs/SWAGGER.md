# Swagger API 文檔使用指南

## 概述

本項目使用 Swagger (OpenAPI) 自動生成 API 文檔，提供交互式的 API 測試界面。

## 訪問 Swagger UI

啟動服務器後，訪問以下地址：

```
http://localhost:8080/swagger/index.html
```

## 功能特性

### 1. 查看所有 API 端點

Swagger UI 會自動顯示所有已配置的 API 端點，按標籤（Tags）分組：
- **認證** - 註冊和登入相關端點
- **用戶** - 用戶管理相關端點
- **系統** - 系統健康檢查等

### 2. 測試 API

在 Swagger UI 中可以直接測試 API：

1. **無需認證的端點**（如 `/register`, `/login`）：
   - 點擊端點展開詳細信息
   - 點擊 "Try it out" 按鈕
   - 填寫請求參數
   - 點擊 "Execute" 執行請求
   - 查看響應結果

2. **需要認證的端點**（如 `/users`, `/profile`）：
   - 首先點擊右上角的 "Authorize" 按鈕
   - 在彈出的對話框中輸入 JWT token（格式：`Bearer {token}` 或直接輸入 token）
   - 點擊 "Authorize" 確認
   - 然後按照上述步驟測試 API

### 3. 查看 API 規範

Swagger UI 會顯示：
- 請求方法（GET, POST 等）
- 端點路徑
- 請求參數（路徑參數、查詢參數、請求體）
- 響應示例和狀態碼
- 數據模型定義

## 生成 Swagger 文檔

當修改了 API 註釋後，需要重新生成文檔：

```bash
# 安裝 swag 工具（如果尚未安裝）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文檔
swag init
```

生成的文檔位於 `docs/` 目錄：
- `docs/swagger.json` - JSON 格式的 API 規範
- `docs/swagger.yaml` - YAML 格式的 API 規範
- `docs/docs.go` - Go 代碼形式的文檔定義

## 使用示例

### 1. 註冊新用戶

在 Swagger UI 中：
1. 找到 `/register` 端點（認證標籤下）
2. 點擊 "Try it out"
3. 修改請求體：
   ```json
   {
     "name": "測試用戶",
     "email": "test@example.com",
     "password": "password123"
   }
   ```
4. 點擊 "Execute"
5. 複製返回的 `token` 值

### 2. 使用 Token 訪問受保護的 API

1. 點擊右上角的 "Authorize" 按鈕
2. 將剛才獲取的 token 粘貼到輸入框（可以直接輸入 token，不需要 "Bearer" 前綴）
3. 點擊 "Authorize"
4. 現在可以測試需要認證的端點（如 `/profile`, `/users`）

### 3. 查看 API 定義

每個端點都會顯示：
- **Summary**: 端點簡短描述
- **Description**: 詳細說明
- **Parameters**: 請求參數列表
- **Responses**: 可能的響應狀態碼和示例

## Swagger 註釋格式

在代碼中添加 Swagger 註釋來生成文檔：

```go
// @Summary 端點簡短描述
// @Description 端點詳細描述
// @Tags 標籤名稱
// @Accept json
// @Produce json
// @Param param_name path/query/body type required "參數說明" example(example_value)
// @Success 200 {object} ResponseType "成功響應說明"
// @Failure 400 {object} map[string]string "錯誤響應說明"
// @Security BearerAuth
// @Router /endpoint [method]
func Handler(c *gin.Context) {
    // ...
}
```

## 常見問題

### Q: Swagger UI 顯示空白頁面？

A: 確保已經運行 `swag init` 生成了文檔，並且 `main.go` 中導入了 `_ "member_API/docs"`。

### Q: 如何更新 Swagger 文檔？

A: 修改代碼中的 Swagger 註釋後，運行 `swag init` 重新生成文檔，然後重啟服務器。

### Q: 如何在 Swagger UI 中使用 JWT token？

A: 點擊右上角的 "Authorize" 按鈕，輸入 token。Swagger UI 會自動在後續請求中添加 `Authorization: Bearer {token}` header。

### Q: 可以導出 Swagger 文檔嗎？

A: 可以，訪問 `http://localhost:8080/swagger/doc.json` 可以下載 JSON 格式的 API 規範。
