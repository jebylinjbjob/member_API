# Dockerfile for Member API
# Build with: docker build -t member-api .

FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o member_api ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates=20241010-r0
WORKDIR /app

COPY --from=builder /src/member_api ./member_api

EXPOSE 9876

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["./member_api"]
