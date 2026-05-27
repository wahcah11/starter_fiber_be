FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git curl

RUN go install github.com/air-verse/air@v1.61.5 && \
    go install golang.org/x/tools/gopls@v0.17.1 && \
    go install golang.org/x/tools/cmd/goimports@v0.30.0

EXPOSE 8080