# 認證功能使用指南

## 概述

本項目實現了完整的 JWT 認證機制，包括用戶註冊、登入和受保護的 API 端點。

## 數據庫設置

首先，確保數據庫表結構包含 `password_hash` 字段：

```sql
-- 如果表已存在，添加字段
ALTER TABLE members ADD COLUMN password_hash VARCHAR(255);

-- 或者執行完整的遷移腳本
psql -d member_api -f database/migration.sql
```

## 環境變量

設置 JWT 密鑰（生產環境必須更改）：

```bash
export JWT_SECRET="your-very-secret-key-change-in-production"
```

## API 使用示例

### 1. 用戶註冊

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "張三",
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

**響應：**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "張三",
    "email": "zhangsan@example.com"
  }
}
```

### 2. 用戶登入

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

**響應：**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "張三",
    "email": "zhangsan@example.com"
  }
}
```

### 3. 訪問受保護的 API

使用登入獲得的 token 訪問受保護的端點：

```bash
# 獲取當前用戶信息
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 獲取所有會員
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 獲取單個會員
curl -X GET http://localhost:8080/api/v1/user/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4. 測試未認證訪問

嘗試不帶 token 訪問受保護的端點：

```bash
curl -X GET http://localhost:8080/api/v1/users
```

**響應：**
```json
{
  "error": "缺少 Authorization header"
}
```

## 公開端點（無需認證）

以下端點不需要認證：

- `POST /api/v1/register` - 用戶註冊
- `POST /api/v1/login` - 用戶登入
- `GET /health` - 健康檢查
- `GET /Hello` - 測試端點
- `GET /graphql` - GraphQL Playground（可選擇是否需要認證）

## 受保護的端點（需要認證）

以下端點需要 JWT token：

- `GET /api/v1/users` - 獲取所有會員
- `GET /api/v1/user/:id` - 獲取單個會員
- `GET /api/v1/profile` - 獲取當前用戶信息

## 安全注意事項

1. **JWT Secret**: 生產環境務必設置強隨機的 `JWT_SECRET`
2. **密碼強度**: 目前要求最少 6 個字符，建議增強驗證規則
3. **HTTPS**: 生產環境必須使用 HTTPS 傳輸 token
4. **Token 過期**: 當前 token 有效期為 24 小時
5. **密碼存儲**: 使用 bcrypt 進行密碼加密，不存儲明文

## 故障排除

### Token 無效錯誤

- 檢查 token 是否過期（24小時）
- 確認 Authorization header 格式正確：`Bearer {token}`
- 確認 `JWT_SECRET` 環境變量與生成 token 時一致

### 密碼驗證失敗

- 確認密碼正確
- 檢查數據庫中 `password_hash` 字段是否正確存儲

### 用戶不存在錯誤

- 確認用戶已註冊
- 檢查數據庫中是否有對應的用戶記錄
