# Gunakan versi stable, misal 1.23
FROM golang:1.25-alpine

# Instal 'air'
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Instal git (penting untuk go mod tidy) dan tools dasar
RUN apk add --no-cache git curl

# Expose port
EXPOSE 8080