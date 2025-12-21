# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o member_api ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /src/member_api ./member_api

EXPOSE 8080
ENV POSTGRES_DSN=postgres://postgres:postgres@db:5432/member_api?sslmode=disable

ENTRYPOINT ["./member_api"]
