# Member API

一個使用 Go、Gin 框架和 PostgreSQL 構建的 RESTful 和 GraphQL API 服務，提供會員管理功能。

## 📋 目錄

-   [技術棧](#技術棧)
-   [功能特性](#功能特性)
-   [項目結構](#項目結構)
-   [快速開始](#快速開始)
-   [環境變量](#環境變量)
-   [API 文檔](#api-文檔)
-   [API 端點](#api-端點)
-   [數據庫設置](#數據庫設置)
-   [資料庫遷移](#資料庫遷移)
-   [Docker 部署](#docker-部署)
-   [開發](#開發)

## 🛠 技術棧

-   **語言**: Go 1.24.x
-   **Web 框架**: Gin
-   **數據庫**: PostgreSQL
-   **API 類型**: RESTful API 和 GraphQL API
-   **API 文檔**: Swagger
-   **容器化**: Docker

## ✨ 功能特性

-   ✅ RESTful API 設計
-   ✅ GraphQL API 支持
-   ✅ PostgreSQL 數據庫集成
-   ✅ Swagger API 文檔
-   ✅ GraphQL Playground
-   ✅ JWT Token 認證
-   ✅ 用戶註冊和登入
-   ✅ 密碼加密（bcrypt）
-   ✅ 健康檢查端點
-   ✅ Docker 容器化支持
-   ✅ 會員管理功能（CRUD）

## 📁 項目結構

```
member_API/
├── main.go                 # 應用程序入口
├── auth/                   # 認證模組
│   ├── jwt.go             # JWT token 處理
│   ├── middleware.go      # 認證中間件
│   └── password.go        # 密碼加密和驗證
├── controllers/            # 控制器層
│   ├── user_controller.go  # 用戶控制器（REST API）
│   └── auth_controller.go  # 認證控制器（登入/註冊）
├── graphql/                # GraphQL 層
│   ├── schema.go          # GraphQL schema 定義
│   ├── resolver.go        # GraphQL resolvers
│   ├── handler.go         # GraphQL handler
│   └── examples.md        # GraphQL 查詢示例
├── routes/                 # 路由層
│   └── routes.go          # 路由配置
├── database/               # 數據庫相關
│   └── migration.sql      # 數據庫遷移腳本
├── go.mod                 # Go 模組依賴
├── go.sum                 # 依賴校驗和
├── Dockerfile             # Docker 構建文件
└── README.md              # 項目文檔
```

## 🚀 快速開始

### 前置要求

-   Go 1.24 或更高版本
-   PostgreSQL 數據庫
-   Git

### 安裝步驟

1. **克隆項目**

    ```bash
    git clone <repository-url>
    cd member_API
    ```

2. **安裝依賴**

    ```bash
    go mod download
    ```

3. **設置環境變量**（可選）

    ```bash
    export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/member_api?sslmode=disable"
    ```

4. **運行應用**
    ```bash
    go run main.go
    ```

服務器將在 `http://localhost:8080` 啟動。

## 🔧 環境變量

| 變量名         | 說明                  | 默認值                                                                   |
| -------------- | --------------------- | ------------------------------------------------------------------------ |
| `POSTGRES_DSN` | PostgreSQL 連接字符串 | `postgres://postgres:postgres@localhost:5432/member_api?sslmode=disable` |
| `JWT_SECRET`   | JWT 簽名密鑰          | `your-secret-key-change-in-production`                                   |

### 環境變量格式

```
postgres://username:password@host:port/database?sslmode=disable
```

## 📚 API 文檔

### Swagger 文檔（REST API）

應用啟動後，可以通過以下地址訪問 Swagger API 文檔：

**推薦：自定義 Swagger UI（帶登入表單）**

```
http://localhost:8080/swagger-ui
```

這個頁面提供了類似 ABP 框架的登入體驗：

-   在頁面頂部直接輸入帳號和密碼登入
-   登入成功後自動設置 JWT token
-   無需手動複製 token 到 Authorization header
-   Token 會自動保存，刷新頁面後自動載入

**原有 Swagger UI（保留兼容性）**

```
http://localhost:8080/swagger/index.html
```

### GraphQL Playground

GraphQL API 提供交互式 Playground，可以通過以下地址訪問：

```
http://localhost:8080/graphql
```

在 Playground 中，您可以：

-   編寫和測試 GraphQL 查詢
-   查看完整的 Schema 定義
-   執行 Mutation 操作
-   查看查詢歷史和結果

## 🔌 API 端點

### 健康檢查

-   **GET** `/health`
    -   檢查服務器狀態和數據庫連接狀態
    -   響應示例：
        ```json
        {
            "status": "OK",
            "postgres_status": "Connected"
        }
        ```

### 認證 API

基礎路徑: `/api/v1`

#### 用戶註冊

-   **POST** `/api/v1/register`
    -   註冊新用戶
    -   請求體：
        ```json
        {
            "name": "張三",
            "email": "zhangsan@example.com",
            "password": "password123"
        }
        ```
    -   響應示例：
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

#### 用戶登入

-   **POST** `/api/v1/login`
    -   用戶登入，獲取 JWT token
    -   請求體：
        ```json
        {
            "email": "zhangsan@example.com",
            "password": "password123"
        }
        ```
    -   響應示例：
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

### 會員 API（需要認證）

基礎路徑: `/api/v1`

**注意：** 以下端點需要認證，請求時需要在 Header 中添加：

```
Authorization: Bearer {your-jwt-token}
```

#### 獲取當前用戶信息

-   **GET** `/api/v1/profile`
    -   獲取當前登入用戶的信息
    -   需要認證：是
    -   響應示例：
        ```json
        {
            "user": {
                "id": 1,
                "name": "張三",
                "email": "zhangsan@example.com"
            }
        }
        ```

#### 獲取所有會員

-   **GET** `/api/v1/users`
    -   獲取會員列表（最多 50 條）
    -   需要認證：是
    -   響應示例：
        ```json
        {
            "users": [
                {
                    "id": 1,
                    "name": "John Doe",
                    "email": "john@example.com"
                }
            ]
        }
        ```

#### 根據 ID 獲取會員

-   **GET** `/api/v1/user/:id`
    -   根據 ID 獲取單個會員信息
    -   需要認證：是
    -   路徑參數：
        -   `id` (必填): 會員 ID
    -   響應示例：
        ```json
        {
            "user": {
                "id": 1,
                "name": "John Doe",
                "email": "john@example.com"
            }
        }
        ```
    -   錯誤響應：
        -   `404`: 會員不存在
        -   `500`: 服務器錯誤

### GraphQL API

基礎路徑: `/graphql`

GraphQL API 提供更靈活的數據查詢和操作方式。訪問 `http://localhost:8080/graphql` 使用 GraphQL Playground。

#### 查詢（Query）

**獲取所有會員**

```graphql
query {
    users {
        id
        name
        email
    }
}
```

**根據 ID 獲取會員**

```graphql
query {
    user(id: 1) {
        id
        name
        email
    }
}
```

#### 變更（Mutation）

**創建新會員**

```graphql
mutation {
    createUser(name: "張三", email: "zhangsan@example.com") {
        id
        name
        email
    }
}
```

**更新會員信息**

```graphql
mutation {
    updateUser(id: 1, name: "李四", email: "lisi@example.com") {
        id
        name
        email
    }
}
```

**刪除會員**

```graphql
mutation {
    deleteUser(id: 1)
}
```

更多示例請參見 `graphql/examples.md` 文件。

### 其他端點

-   **GET** `/Hello`
    -   測試端點
    -   響應示例：
        ```json
        {
            "message": "Hello, RESTful API!"
        }
        ```

## 🗄 數據庫設置

### PostgreSQL 數據庫結構

確保數據庫中存在 `members` 表，表結構如下：

```sql
CREATE TABLE members (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**注意：** 如果您的數據庫中已經有 `members` 表，需要添加 `password_hash` 字段：

```sql
ALTER TABLE members ADD COLUMN password_hash VARCHAR(255);
```

### 創建數據庫

```sql
CREATE DATABASE member_api;
```

### 連接數據庫

應用程序會自動連接到配置的 PostgreSQL 數據庫。如果連接失敗，應用程序仍會啟動，但會顯示警告信息。

## 🧱 資料庫遷移

本專案使用 [golang-migrate](https://github.com/golang-migrate/migrate) 內嵌遷移，每次應用啟動時都會自動比對 `migrations` 目錄中的 SQL 並執行 `up` 動作。這確保不同環境中的資料表結構保持一致。

### 目錄結構

```
migrations/
├── 0001_create_members_table.up.sql
└── 0001_create_members_table.down.sql
```

### 新增遷移檔

1. 複製現有檔名格式 `YYYY_description.{up,down}.sql`（數字遞增即可，不一定要年份）。
2. 在 `.up.sql` 中撰寫升級 SQL，在 `.down.sql` 中加入相對應的回滾 SQL。
3. 提交變更後，應用程式會自動套用。

### 手動執行遷移（可選）

若需要單獨執行遷移，可安裝 CLI：

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
```

再透過以下指令控制：

```bash
migrate -path migrations -database "$POSTGRES_DSN" up    # 套用
migrate -path migrations -database "$POSTGRES_DSN" down  # 回滾
```

## 🐳 Docker 部署

### 構建 Docker 鏡像

```bash
docker build -t member-api .
```

### 運行容器

```bash
docker run -p 8080:8080 \
  -e POSTGRES_DSN="postgres://postgres:postgres@db:5432/member_api?sslmode=disable" \
  member-api
```

### 使用 Docker Compose（推薦）

創建 `docker-compose.yml` 文件：

```yaml
version: "3.8"

services:
    db:
        image: postgres:15
        environment:
            POSTGRES_USER: postgres
            POSTGRES_PASSWORD: postgres
            POSTGRES_DB: member_api
        ports:
            - "5432:5432"
        volumes:
            - postgres_data:/var/lib/postgresql/data

    api:
        build: .
        ports:
            - "8080:8080"
        environment:
            POSTGRES_DSN: postgres://postgres:postgres@db:5432/member_api?sslmode=disable
        depends_on:
            - db

volumes:
    postgres_data:
```

運行服務：

```bash
docker-compose up -d
```

## 💻 開發

### 構建應用

```bash
go build -o member_API.exe main.go
```

### 運行測試

```bash
go test ./...
```

### 代碼檢查

```bash
go vet ./...
```

### 格式化代碼

```bash
go fmt ./...
```

### 生成 Swagger 文檔

```bash
# 安裝 swag 工具（如果尚未安裝）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文檔
swag init

# 修復可能的編譯錯誤（Windows PowerShell）
.\scripts\fix-swagger.ps1

# 或在 Linux/Mac 上
chmod +x scripts/fix-swagger.sh
./scripts/fix-swagger.sh

# 文檔將生成在 docs/ 目錄下
```

**注意：**

-   每次修改 API 註釋後，需要重新運行 `swag init` 來更新文檔
-   如果遇到 `LeftDelim` 和 `RightDelim` 編譯錯誤，運行修復腳本即可自動修復

## 📝 注意事項

-   應用程序需要 PostgreSQL 數據庫支持
-   默認端口為 8080
-   數據庫連接使用環境變量 `POSTGRES_DSN` 配置
-   JWT 簽名密鑰使用環境變量 `JWT_SECRET` 配置（生產環境請務必更改）
-   Swagger 文檔路徑為 `/swagger/index.html`
-   GraphQL Playground 路徑為 `/graphql`
-   同時支持 RESTful API 和 GraphQL API
-   大部分 API 端點需要 JWT 認證，請先註冊/登入獲取 token

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 許可證

[在此添加許可證信息]

---
