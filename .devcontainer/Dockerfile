FROM golang:1.24-alpine

WORKDIR /app

# Tambah make, bash, dan tools lain yang dibutuhkan
RUN apk add --no-cache \
    git \
    curl \
    make \
    bash \
    gcc \
    musl-dev

# Go tools
RUN go install github.com/air-verse/air@v1.61.5 && \
    go install golang.org/x/tools/gopls@v0.17.1 && \
    go install golang.org/x/tools/cmd/goimports@v0.30.0 && \
    go install github.com/vektra/mockery/v2@latest

EXPOSE 8080