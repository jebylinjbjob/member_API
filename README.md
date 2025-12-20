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
-   **ORM / 遷移**: GORM（AutoMigrate）
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

** Swagger UI（保留兼容性）**

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

## 🧱 資料庫遷移

目前改用 [GORM](https://gorm.io/) 的 `AutoMigrate` 管理 schema。應用啟動時會自動對 `models.Member` 等模型執行遷移，確保必需的資料表與欄位存在，無須再維護外部 SQL 腳本。

### 新增或修改欄位

1. 在 `models/` 目錄中更新或新增模型結構（例如 `Member`）。
2. 重建或重新啟動服務，`AutoMigrate` 會自動同步結構。

> `AutoMigrate` 只會新增欄位/索引，不會刪除資料或危險操作；若需要更複雜的遷移，建議以 GORM callbacks 或額外腳本實作。

### 可選：手動 SQL

仍可在 `migrations/` 目錄放置輔助 SQL，供 DBA 或 CI/CD 參考，但系統執行時不再自動讀取該目錄。如需手動執行特定 SQL，可自行使用 `psql` 或其他工具。

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
