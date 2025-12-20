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
-   [Docker 部署](#docker-部署)
-   [開發](#開發)

## 🛠 技術棧

-   **語言**: Go 1.23.4
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
-   ✅ 健康檢查端點
-   ✅ Docker 容器化支持
-   ✅ JWT Token 認證（規劃中）
-   ✅ 會員管理功能（CRUD）

## 📁 項目結構

```
member_API/
├── main.go                 # 應用程序入口
├── controllers/            # 控制器層
│   └── user_controller.go  # 用戶控制器（REST API）
├── graphql/                # GraphQL 層
│   ├── schema.go          # GraphQL schema 定義
│   ├── resolver.go        # GraphQL resolvers
│   ├── handler.go         # GraphQL handler
│   └── examples.md        # GraphQL 查詢示例
├── routes/                 # 路由層
│   └── routes.go          # 路由配置
├── go.mod                 # Go 模組依賴
├── go.sum                 # 依賴校驗和
├── Dockerfile             # Docker 構建文件
└── README.md              # 項目文檔
```

## 🚀 快速開始

### 前置要求

-   Go 1.23.4 或更高版本
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

### 環境變量格式

```
postgres://username:password@host:port/database?sslmode=disable
```

## 📚 API 文檔

### Swagger 文檔（REST API）

應用啟動後，可以通過以下地址訪問 Swagger API 文檔：

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

### 會員 API

基礎路徑: `/api/v1`

#### 獲取所有會員

-   **GET** `/api/v1/users`
    -   獲取會員列表（最多 50 條）
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 創建數據庫

```sql
CREATE DATABASE member_api;
```

### 連接數據庫

應用程序會自動連接到配置的 PostgreSQL 數據庫。如果連接失敗，應用程序仍會啟動，但會顯示警告信息。

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

## 📝 注意事項

-   應用程序需要 PostgreSQL 數據庫支持
-   默認端口為 8080
-   數據庫連接使用環境變量 `POSTGRES_DSN` 配置
-   Swagger 文檔路徑為 `/swagger/index.html`
-   GraphQL Playground 路徑為 `/graphql`
-   同時支持 RESTful API 和 GraphQL API

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 許可證

[在此添加許可證信息]

---

**開發中功能**:

-   🔄 JWT Token 認證
