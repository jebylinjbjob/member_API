default:
    @just --list

vet:
    echo "Running go vet..."
    go vet ./...

test:
    echo "Running tests..."
    go test -v ./...

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
