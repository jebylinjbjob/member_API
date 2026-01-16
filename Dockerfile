# Dockerfile for Member API
# Build with: docker build -t member-api .

FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o member_api ./main.go

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /src/member_api ./member_api

EXPOSE 9876

USER 65532:65532

ENTRYPOINT ["./member_api"]
