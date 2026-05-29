# =========================================
# VARIABLES
# =========================================
APP_NAME=starter-wahcah-be
MAIN_PATH=cmd/api/main.go
BINARY_PATH=bin/$(APP_NAME)
MOCK_DIR=internal/modules/auth/login/mocks
MOCK_SRC=internal/modules/auth/login


# =========================================
# PROJECT INIT
# =========================================

## Inisialisasi nama project baru (jalankan sekali setelah clone)
## Contoh: make init PROJECT=my-awesome-api
init:
	@if [ -z "$(PROJECT)" ]; then \
		sh init-project.sh; \
	else \
		sh init-project.sh $(PROJECT); \
	fi


# =========================================
# DEVELOPMENT
# =========================================

## Jalankan server langsung
run:
	go run $(MAIN_PATH)

## Jalankan dengan live reload menggunakan air
dev:
	air

## Build binary ke folder bin/
build:
	go build -o $(BINARY_PATH) $(MAIN_PATH)

## Jalankan binary hasil build
start:
	./$(BINARY_PATH)


# =========================================
# DATABASE & SEEDER
# =========================================

## Jalankan migrate + seeder (otomatis jalan saat run di APP_ENV != production)
seed:
	APP_ENV=development go run $(MAIN_PATH)

## Tidy dependencies
tidy:
	go mod tidy


# =========================================
# MOCK & TEST
# =========================================

## Generate mock dari semua interface di internal/modules
mock:
	mockery

## Jalankan semua test (unit + integration + util + middleware)
test:
	APP_ENV=test go test ./internal/... -v

## Hanya unit test modules (tanpa integration)
test-unit:
	APP_ENV=test go test ./internal/modules/... -v

## Hanya util dan middleware test
test-util:
	APP_ENV=test go test ./internal/util/... ./internal/middleware/... -v

## Hanya integration test
test-integration:
	APP_ENV=test go test ./internal/modules/... -v -run Integration

## Jalankan test dengan coverage report (HTML)
test-coverage:
	APP_ENV=test go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Jalankan test satu package (contoh: make test-pkg PKG=./internal/modules/auth/login/...)
test-pkg:
	APP_ENV=test go test $(PKG) -v


# =========================================
# CODE QUALITY
# =========================================

## Format semua file Go
fmt:
	go fmt ./...

## Jalankan linter (butuh golangci-lint)
lint:
	golangci-lint run ./...

## Format + lint sekaligus
check: fmt lint


# =========================================
# UTILITY
# =========================================

## Hapus binary dan coverage output
clean:
	rm -rf bin/ coverage.out coverage.html

## Tampilkan semua perintah yang tersedia
help:
	@echo ""
	@echo "========================================"
	@echo "  $(APP_NAME)"
	@echo "========================================"
	@echo ""
	@echo "Project Init:"
	@echo "  make init                    Inisialisasi nama project (interaktif)"
	@echo "  make init PROJECT=nama-baru  Inisialisasi nama project langsung"
	@echo ""
	@echo "Development:"
	@echo "  make run                     Jalankan server"
	@echo "  make dev                     Jalankan dengan live reload (air)"
	@echo "  make build                   Build binary ke bin/"
	@echo "  make start                   Jalankan binary hasil build"
	@echo ""
	@echo "Database:"
	@echo "  make seed                    Migrate + seeder (development)"
	@echo "  make tidy                    Tidy dependencies"
	@echo ""
	@echo "Mock & Test:"
	@echo "  make mock                    Generate semua mock di internal/modules"
	@echo "  make test                    Semua test (unit + integration + util + middleware)"
	@echo "  make test-unit               Hanya test modules"
	@echo "  make test-util               Hanya test util dan middleware"
	@echo "  make test-integration        Hanya integration test"
	@echo "  make test-coverage           Test + coverage report HTML"
	@echo "  make test-pkg PKG=./path/... Test satu package spesifik"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt                     Format kode"
	@echo "  make lint                    Jalankan linter"
	@echo "  make check                   fmt + lint sekaligus"
	@echo ""
	@echo "Utility:"
	@echo "  make clean                   Hapus bin/ dan coverage output"
	@echo "  make help                    Tampilkan bantuan ini"
	@echo ""

.PHONY: init run dev build start seed tidy mock test test-unit test-util test-integration test-coverage test-pkg fmt lint check clean help