# 會員 API 專案的 Justfile
default:
    @just --list

# 本地運行
run:
    go run main.go

# 編譯
build:
    go build -o bin/member_api main.go

# 測試
test:
    go test -v ./...

# 格式化
fmt:
    go fmt ./...

# 整理依賴
tidy:
    go mod tidy

# 測試覆蓋率
test-cov:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html


# 生成 GraphQL
graphql:
    go get github.com/99designs/gqlgen@v0.17.85
    go run github.com/99designs/gqlgen generate

# 清理
clean:
    rm -rf bin/ coverage.out coverage.html

# 一鍵檢查
auto_check:
    go mod tidy
    go fmt ./...
    go test -v ./...
    go build
    just test-cov