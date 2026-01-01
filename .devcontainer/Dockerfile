# Gunakan versi Go stable
FROM golang:1.25-alpine

# Set environment PATH supaya 'air' bisa dipanggil
ENV PATH=$PATH:/go/bin

# Install tools dasar: bash, mysql client, git, curl
RUN apk add --no-cache bash mysql-client git curl

# Install Air (hot reload)
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go.mod + go.sum dulu untuk caching layer
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy semua source code
COPY . .

# Expose port backend
EXPOSE 8080

# Default command: jalankan Air
CMD ["air"]
