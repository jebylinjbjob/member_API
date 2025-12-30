default:
    @just --list

vet:
    echo "Running go vet..."
    go vet ./...

# 單元測試（使用 Mock，快速，不帶 build tag）
test:
    echo "Running unit tests (without integration tag)..."
    go test -v ./...

# 整合測試（使用 SQLite，較慢，需要 integration build tag）
# 設置 CGO_ENABLED=0 以使用純 Go 的 SQLite 驅動（不需要 CGO）
test-integration:
    echo "Running integration tests..."
    CGO_ENABLED=0 go test -v -tags=integration ./...

# 運行所有測試（單元 + 整合）
test-all:
    echo "Running all tests..."
    just test
    just test-integration

install-tools:
    echo "Installing development tools..."
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build:
    echo "Building application..."
    go build -o member_API.exe main.go

run:
    echo "Running application..."
    go run main.go
