# Member API

一個使用 Go、Gin 框架和 PostgreSQL 構建的 RESTful 和 GraphQL API 服務，提供會員管理功能和 JWT 認證。

## 功能特性

- ✅ RESTful API 端點
- ✅ GraphQL API 支援
- ✅ JWT 認證機制
- ✅ 用戶註冊與登入
- ✅ 會員管理功能
- ✅ 產品管理功能（CRUD）
- ✅ 密碼加密（bcrypt）
- ✅ Repository 模式架構
- ✅ 單元測試與整合測試
- ✅ PostgreSQL 資料庫
- ✅ Swagger API 文件
- ✅ Docker 支援
- ✅ 自動更新部署（Watchtower）
- ✅ CI/CD 自動化測試

## 技術棧

- **語言**: Go 1.24+
- **框架**: Gin Web Framework
- **資料庫**: PostgreSQL（生產環境）、SQLite（測試環境）
- **ORM**: GORM
- **認證**: JWT (golang-jwt/jwt)
- **密碼加密**: bcrypt
- **API 文件**: Swagger/OpenAPI
- **GraphQL**: gqlgen
- **測試框架**: testify/mock（單元測試）
- **開發工具**: Just（命令執行器）
- **容器化**: Docker & Docker Compose
- **自動更新**: Watchtower
- **CI/CD**: GitHub Actions

## 快速開始

### 方式一：使用 Docker Compose（推薦）

這是最簡單的部署方式，包含自動更新功能。

1. 確保已安裝 Docker 和 Docker Compose
2. 編輯 `deploy/docker-compose.yml` 設定資料庫連線和 JWT 密鑰
3. 啟動服務：

```bash
docker-compose -f deploy/docker-compose.yml up -d
```

4. 查看日誌：

```bash
docker-compose -f deploy/docker-compose.yml logs -f member_api
```

5. 停止服務：

```bash
docker-compose -f deploy/docker-compose.yml down
```

**自動更新機制**：

- Watchtower 每 5 分鐘檢查一次新版本
- 發現新版本時自動拉取並重啟容器
- 自動清理舊映像檔
- 無需手動干預

### 方式二：本地開發

#### 環境需求

- Go 1.24 或更高版本
- PostgreSQL 12 或更高版本
- Git

#### 環境變數設定

複製環境變數範例檔案並根據需要修改：

```bash
cp .env.example .env
```

編輯 `.env` 檔案設定以下變數：

- `POSTGRES_DSN`: PostgreSQL 資料庫連線字串
- `JWT_SECRET`: JWT 簽章密鑰（**生產環境務必更改**）

生成強隨機密鑰的方法：

```bash
# 使用 openssl 生成 32 位元組的 base64 編碼密鑰
openssl rand -base64 32
```

#### 安裝相依套件

```bash
go mod download
```

#### 安裝開發工具（可選）

專案使用 [Just](https://github.com/casey/just) 作為命令執行器，提供便捷的開發命令：

```bash
# 安裝 Just（如果尚未安裝）
# Windows (使用 Scoop)
scoop install just

# macOS (使用 Homebrew)
brew install just

# Linux
# 參考 https://github.com/casey/just#installation

# 查看所有可用命令
just
```

#### 生成 Swagger 文件

```bash
# 安裝 swag（如果尚未安裝）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文件
swag init
```

#### 啟動服務

使用 Just（推薦）：

```bash
just run
```

服務將在 `http://localhost:9876` 啟動。

## Docker 部署

### 從 GitHub Container Registry 抓取映像檔

```bash
# 抓取最新版本
docker pull ghcr.io/jebylinjbjob/member_api:latest

# 或抓取特定版本
docker pull ghcr.io/jebylinjbjob/member_api:v1.0.0
```

### 本地建構映像檔

```bash
docker build -t member-api .
```

### 執行容器

```bash
docker run -d \
  -p 9876:9876 \
  -e POSTGRES_DSN="postgres://postgres:postgres@host.docker.internal:5432/member_api?sslmode=disable" \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  --name member-api \
  ghcr.io/jebylinjbjob/member_api:latest
```

## API 文件

### Swagger UI

啟動服務後，訪問 Swagger UI：

```
http://localhost:9876/swagger/index.html
```

### GraphQL Playground

訪問 GraphQL Playground：

```
http://localhost:9876/graphql
```

### 健康檢查

```bash
curl http://localhost:9876/health
```

## API 端點

### 公開端點（無需認證）

- `POST /api/v1/register` - 用戶註冊
- `POST /api/v1/login` - 用戶登入
- `GET /health` - 健康檢查

### 受保護端點（需要 JWT 認證）

#### 會員管理

- `GET /api/v1/users` - 獲取所有會員
- `GET /api/v1/user/:id` - 獲取單個會員
- `GET /api/v1/profile` - 獲取當前用戶資訊
- `DELETE /api/v1/user/:id` - 刪除會員（軟刪除）

#### 產品管理

- `GET /api/v1/products` - 獲取產品列表（支援分頁）
  - Query 參數：`limit`（預設 50，最大 100）、`offset`（預設 0）
- `GET /api/v1/product/:id` - 獲取單個產品
- `POST /api/v1/product` - 創建產品
- `PUT /api/v1/product/:id` - 更新產品
- `DELETE /api/v1/product/:id` - 刪除產品（軟刪除）

### 使用 JWT 認證

1. 先登入獲取 token：

```bash
curl -X POST http://localhost:9876/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

2. 使用 token 訪問受保護的端點：

```bash
curl -X GET http://localhost:9876/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

3. 產品管理範例：

```bash
# 獲取產品列表（分頁）
curl -X GET "http://localhost:9876/api/v1/products?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 創建產品
curl -X POST http://localhost:9876/api/v1/product \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15 Pro",
    "product_price": 35900,
    "product_description": "最新款 iPhone",
    "product_image": "https://example.com/image.jpg",
    "product_stock": 100
  }'

# 更新產品
curl -X PUT http://localhost:9876/api/v1/product/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15 Pro Max",
    "product_price": 42900
  }'
```

## 專案結構

```
member_API/
├── auth/              # JWT 認證與密碼處理
│   ├── jwt.go
│   ├── middleware.go
│   └── password.go
├── controllers/       # 控制器層（處理 HTTP 請求）
│   ├── auth_controller.go
│   ├── product_controller.go
│   └── user_controller.go
├── repositories/     # Repository 層（資料存取抽象）
│   ├── gorm_member_repository.go
│   ├── gorm_product_repository.go
│   ├── member_repository.go
│   ├── product_repository.go
│   ├── mock_member_repository.go
│   └── mock_product_repository.go
├── services/          # Service 層（業務邏輯）
│   ├── member_service.go
│   ├── member_service_test.go
│   ├── member_service_integration_test.go
│   ├── product_service.go
│   ├── product_service_test.go
│   └── product_service_integration_test.go
├── models/            # 資料庫結構
│   ├── base.go        #資料庫通用結構
│   ├── member.go
│   └── product.go
├── test/              # 測試輔助工具
│   └── test_helpers.go
├── deploy/            # 部署配置
│   └── docker-compose.yml
├── docs/              # API 文件
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── graphql/           # GraphQL schema 和 resolver
│   ├── generated.go
│   ├── resolver.go
│   ├── schema.graphql
│   ├── schema.resolvers.go
│   ├── setup.go
│   └── model/
│       └── models_gen.go
├── routes/            # 路由設定
│   └── routes.go
├── .github/            # GitHub Actions 工作流程
│   └── workflows/
│       ├── LintTest.yaml
│       └── Trivy.yml
├── Dockerfile         # Docker 建構檔案
├── justfile           # Just 命令定義
├── .editorconfig      # 編輯器配置
├── go.mod             # Go 模組定義
├── gqlgen.yml         # GraphQL 程式碼產生器配置
├── main.go            # 應用程式入口
└── .env.example       # 環境變數範例
```

### 架構說明

專案採用**分層架構**和**Repository 模式**：

- **Controllers**: 處理 HTTP 請求/響應，驗證輸入
- **Services**: 實現業務邏輯，不直接依賴資料庫
- **Repositories**: 抽象資料存取層，提供介面定義和 GORM 實現
- **Models**: 定義資料結構

這種設計的優點：

- 關注點分離，易於維護
- 易於單元測試（使用 Mock Repository）
- 可替換資料庫實現（只需實現 Repository 介面）

## 開發指南

### 使用 Just 命令（推薦）

專案提供了便捷的 Just 命令來執行常見開發任務：

```bash
# 查看所有可用命令
just

# 執行單元測試（使用 Mock，快速）
just test

# [開發中] 執行整合測試（使用 SQLite，較慢）
just test-integration

# 執行所有測試（單元 + 整合）
just test-all

# 程式碼檢查
just vet

# 建構應用程式
just build

# 執行應用程式
just run

# 安裝開發工具（goimports, golangci-lint）
just install-tools
```

### 執行測試

專案包含兩種測試類型：

#### 單元測試（使用 Mock）

快速執行，不依賴資料庫：

```bash
# 使用 Just
just test

# 或直接使用 Go
go test -v ./...
```

#### 整合測試（使用 SQLite）

需要 `integration` build tag，使用真實資料庫：

```bash
# 使用 Just
just test-integration

# 或直接使用 Go
CGO_ENABLED=0 go test -v -tags=integration ./...
```

#### 執行所有測試

```bash
just test-all
```

### 程式碼格式化

```bash
go fmt ./...
```

### 程式碼檢查

```bash
# 使用 Just
just vet

# 或直接使用 Go
go vet ./...
```

### 使用 Linter

```bash
# 安裝 golangci-lint（如果尚未安裝）
just install-tools

# 執行 linter
golangci-lint run
```

### 更新 Swagger 文件

當 API 端點有變更時，重新生成 Swagger 文件：

```bash
swag init
```

### GraphQL Schema 變更

編輯 `graphql/schema.graphql` 後，執行以下命令重新生成程式碼：

```bash
go run github.com/99designs/gqlgen generate
```

## 環境變數說明

| 變數名稱       | 說明                      | 範例                                                             |
| -------------- | ------------------------- | ---------------------------------------------------------------- |
| `POSTGRES_DSN` | PostgreSQL 資料庫連線字串 | `postgres://user:password@localhost:5432/dbname?sslmode=disable` |
| `JWT_SECRET`   | JWT 簽章密鑰              | `your-secret-key`                                                |

## 生產環境建議

1. **更改預設密鑰**：務必使用強隨機字串作為 `JWT_SECRET`
2. **啟用 SSL**：在資料庫連線字串中將 `sslmode=disable` 改為 `sslmode=require`
3. **設定防火牆**：僅開放必要的連接埠
4. **定期備份**：設定資料庫自動備份機制
5. **監控日誌**：使用日誌聚合工具監控應用程式狀態
6. **限流**：在生產環境中加入 API 限流機制

## 疑難排解

### 資料庫連線失敗

檢查 `POSTGRES_DSN` 是否正確，確保：

- 資料庫伺服器正在執行
- 連線字串中的使用者名稱、密碼、主機、連接埠、資料庫名稱都正確
- 網路連線正常

### Docker 容器無法啟動

```bash
# 查看容器日誌
docker logs member_api

# 檢查容器狀態
docker ps -a
```

### Watchtower 未自動更新

```bash
# 查看 Watchtower 日誌
docker logs member_api_watchtower

# 手動觸發更新
docker-compose -f deploy/docker-compose.yml pull
docker-compose -f deploy/docker-compose.yml up -d
```

## CI/CD

專案使用 GitHub Actions 進行自動化測試和程式碼檢查：

- **Lint**: 使用 `golangci-lint` 檢查程式碼品質
- **單元測試**: 自動執行單元測試並上傳覆蓋率報告
- **整合測試**: 自動執行整合測試並上傳覆蓋率報告
- **Docker 檢查**: 使用 `hadolint` 檢查 Dockerfile

所有測試和檢查會在每次 push 和 pull request 時自動執行。

## 程式碼風格

專案使用 `.editorconfig` 確保一致的程式碼風格：

- Go 檔案：使用 Tab 縮排，大小為 4
- YAML/JSON 檔案：使用空格縮排，大小為 2
- 自動移除尾隨空白
- 統一使用 LF 行尾符號

大多數現代編輯器都支援 `.editorconfig`，會自動應用這些設定。
