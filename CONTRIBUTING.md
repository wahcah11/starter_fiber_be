# Panduan Penulisan Kode — starter-wahcah-be

Dokumen ini menjelaskan aturan penulisan kode untuk setiap layer, mulai dari penamaan file, interface, struct, hingga unit test. Wajib dibaca sebelum menambahkan module baru.

---

## Daftar Isi

1. [Memulai Project Baru](#memulai-project-baru)
2. [Struktur Folder](#struktur-folder)
3. [Aturan Umum](#aturan-umum)
4. [Layer DTO](#layer-dto)
5. [Layer Repository](#layer-repository)
6. [Layer Service](#layer-service)
7. [Layer Controller](#layer-controller)
8. [Layer Router](#layer-router)
9. [Unit Test per Layer](#unit-test-per-layer)
10. [Test Util](#test-util)
11. [Test Middleware](#test-middleware)
12. [Integration Test Repository](#integration-test-repository)
13. [Checklist Membuat Module Baru](#checklist-membuat-module-baru)
14. [Perintah Penting](#perintah-penting)

---

## Memulai Project Baru

Starter ini dirancang untuk di-clone lalu langsung diganti namanya menjadi project Anda. Ikuti langkah berikut secara berurutan.

### Langkah 1 — Clone Repository

```bash
git clone https://github.com/your-username/starter-wahcah-be.git nama-project-baru
cd nama-project-baru
```

### Langkah 2 — Inisialisasi Nama Project

Jalankan script init untuk mengganti semua referensi `starter-wahcah-be` menjadi nama project Anda. Script ini akan menyentuh semua file `*.go`, `go.mod`, `.env`, `*.yaml`, `*.yml`, `*.md`, dan `Makefile` sekaligus.

**Interaktif** — akan ditanya nama project:
```bash
sh init-project.sh
```

**Langsung dengan argumen:**
```bash
sh init-project.sh my-awesome-api
# atau
make init PROJECT=my-awesome-api
```

Contoh output:
```
Mengganti 'starter-wahcah-be' -> 'my-awesome-api' ...

  updated: go.mod
  updated: Makefile
  updated: cmd/api/main.go
  updated: config/database.go
  updated: internal/modules/auth/login/login_service.go
  updated: internal/seeder/user_seeder.go
  updated: CONTRIBUTING.md
  ...

Selesai! Project berhasil diinisialisasi sebagai 'my-awesome-api'.
```

> **Catatan:** Script hanya boleh dijalankan **sekali** setelah clone. Jika dijalankan ulang, nama lama sudah tidak ada dan tidak ada yang berubah.

### Langkah 3 — Sesuaikan .env

Buka file `.env` dan sesuaikan konfigurasi:

```env
APP_NAME=my-awesome-api
APP_ENV=development

DB_HOST=db
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

JWT_SECRET=ganti-dengan-string-acak-yang-panjang
```

Generate JWT_SECRET yang aman:
```bash
openssl rand -base64 32
```

### Langkah 4 — Install Dependencies

```bash
make tidy
```

### Langkah 5 — Generate Mock

```bash
make mock
```

### Langkah 6 — Jalankan Test

Pastikan semua test bawaan starter masih PASS sebelum mulai development:

```bash
make test
```

Output yang diharapkan:
```
PASS
ok    my-awesome-api/internal/modules/auth/login   5.4s
ok    my-awesome-api/internal/util                 0.1s
ok    my-awesome-api/internal/middleware            0.1s
```

### Langkah 7 — Jalankan Server

```bash
make dev
```

Server berjalan di `http://localhost:8080`. Health check tersedia di:
```bash
curl http://localhost:8080/health
# {"status":"ok","env":"development"}
```

### Langkah 8 — Reset Git History (Opsional)

Jika tidak ingin membawa history commit starter ini:

```bash
rm -rf .git
git init
git add .
git commit -m "feat: init project my-awesome-api"
```

---


## Struktur Folder

Setiap module berada di dalam `internal/modules/{domain}/{feature}/`.

```
internal/
├── middleware/
│   ├── auth_middleware.go
│   └── auth_middleware_test.go
├── models/
│   ├── migrate.go
│   ├── user_model.go
│   ├── role_model.go
│   ├── permission_model.go
│   └── user_role_model.go
├── modules/
│   └── {domain}/
│       └── {feature}/
│           ├── {feature}_dto.go
│           ├── {feature}_repository.go
│           ├── {feature}_service.go
│           ├── {feature}_controller.go
│           ├── {feature}_router.go
│           ├── {feature}_repository_test.go
│           ├── {feature}_service_test.go
│           ├── {feature}_controller_test.go
│           ├── {feature}_router_test.go
│           ├── {feature}_repository_integration_test.go
│           └── mocks/
│               ├── Repository.go   <- auto-generated mockery
│               └── Service.go      <- auto-generated mockery
├── router/
│   └── router.go
└── util/
    ├── jwt_util.go
    ├── jwt_util_test.go
    ├── password_util.go
    ├── password_util_test.go
    ├── validation_util.go
    └── validation_util_test.go
```

Contoh nyata untuk fitur login:

```
internal/modules/auth/login/
├── login_dto.go
├── login_repository.go
├── login_service.go
├── login_controller.go
├── login_router.go
├── login_repository_test.go
├── login_service_test.go
├── login_controller_test.go
├── login_router_test.go
├── login_repository_integration_test.go
└── mocks/
    ├── Repository.go
    └── Service.go
```

---

## Aturan Umum

### Penamaan File

| Layer | Format Nama File | Contoh |
|---|---|---|
| DTO | `{feature}_dto.go` | `login_dto.go` |
| Repository | `{feature}_repository.go` | `login_repository.go` |
| Service | `{feature}_service.go` | `login_service.go` |
| Controller | `{feature}_controller.go` | `login_controller.go` |
| Router | `{feature}_router.go` | `login_router.go` |
| Test Repository | `{feature}_repository_test.go` | `login_repository_test.go` |
| Test Service | `{feature}_service_test.go` | `login_service_test.go` |
| Test Controller | `{feature}_controller_test.go` | `login_controller_test.go` |
| Test Router | `{feature}_router_test.go` | `login_router_test.go` |
| Integration Test | `{feature}_repository_integration_test.go` | `login_repository_integration_test.go` |

### Penamaan Package

Semua file dalam satu folder feature menggunakan package yang sama, yaitu nama feature-nya.

```go
// Semua file di folder login/ menggunakan:
package login

// Semua file test di folder login/ menggunakan:
package login_test
```

### Penamaan Struct & Interface

| Jenis | Format | Contoh |
|---|---|---|
| Interface Repository | `Repository` | `Repository` |
| Struct Repository | `repository` (huruf kecil) | `repository` |
| Constructor Repository | `New{Feature}Repository` | `NewLoginRepository` |
| Interface Service | `Service` | `Service` |
| Struct Service | `service` (huruf kecil) | `service` |
| Constructor Service | `New{Feature}Service` | `NewLoginService` |
| Struct Controller | `Controller` (huruf besar) | `Controller` |
| Constructor Controller | `New{Feature}Controller` | `NewLoginController` |
| Fungsi Router | `InitRoutes` | `InitRoutes` |

> **Aturan:** Interface selalu huruf besar (exported). Struct implementasi selalu huruf kecil (unexported) — hanya bisa diakses melalui constructor dan interface-nya.

---

## Layer DTO

DTO (Data Transfer Object) adalah struct untuk request dan response HTTP. Tidak boleh ada logic di sini.

### Aturan

- Nama file: `{feature}_dto.go`
- Nama struct Request: `{Feature}Request`
- Nama struct Response: `{Feature}Response`
- Selalu tambahkan tag `json` dan `validate`
- Tidak boleh import model database langsung

### Contoh

```go
package login

// LoginRequest adalah struct untuk request body POST /auth/login
type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse adalah struct untuk response token setelah login berhasil
type LoginResponse struct {
    Token string `json:"token"`
}
```

### Tag Validasi yang Umum Dipakai

| Tag | Keterangan |
|---|---|
| `required` | Field wajib diisi |
| `email` | Harus format email valid |
| `min=N` | Minimum N karakter |
| `max=N` | Maksimum N karakter |
| `oneof=a b c` | Harus salah satu dari nilai tersebut |

---

## Layer Repository

Repository bertanggung jawab **hanya** untuk operasi database. Tidak boleh ada logic bisnis di sini.

### Aturan

- Nama file: `{feature}_repository.go`
- Wajib ada **interface** `Repository` dan **struct** `repository`
- Struct `repository` hanya boleh menyimpan `*gorm.DB`
- Setiap method hanya boleh melakukan satu operasi database
- Kembalikan `error` dari GORM langsung, jangan dibungkus pesan custom
- Tidak boleh import package `service`, `controller`, atau `router`

### Struktur Wajib

```go
package login

import (
    "starter-wahcah-be/internal/models"
    "gorm.io/gorm"
)

// Repository mendefinisikan kontrak operasi database untuk login
type Repository interface {
    FindByEmail(email string) (*models.User, error)
}

// repository adalah implementasi dari interface Repository
type repository struct {
    db *gorm.DB
}

// NewLoginRepository membuat instance repository baru
func NewLoginRepository(db *gorm.DB) Repository {
    return &repository{db}
}

// FindByEmail mencari user berdasarkan email
func (r *repository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

### Yang Boleh dan Tidak Boleh

| Boleh | Tidak Boleh |
|---|---|
| `db.Where(...).First(...)` | Logic if/else bisnis |
| `db.Create(...)` | Hash password |
| `db.Save(...)` | Generate token |
| `db.Delete(...)` | Return pesan error custom |
| Return error GORM langsung | Import service/controller |

---

## Layer Service

Service bertanggung jawab untuk **logic bisnis**. Menggunakan repository melalui interface, tidak langsung ke database.

### Aturan

- Nama file: `{feature}_service.go`
- Wajib ada **interface** `Service` dan **struct** `service`
- Struct `service` hanya boleh menyimpan `Repository` (interface, bukan struct)
- Semua error yang dikembalikan ke controller harus berupa pesan yang ramah pengguna
- Tidak boleh import `fiber`, `gorm`, atau package HTTP lainnya
- Tidak boleh akses database langsung — selalu lewat repository

### Struktur Wajib

```go
package login

import (
    "errors"
    "starter-wahcah-be/internal/util"
)

// Service mendefinisikan kontrak logic bisnis untuk login
type Service interface {
    Authenticate(req LoginRequest) (*LoginResponse, error)
}

// service adalah implementasi dari interface Service
type service struct {
    repo Repository
}

// NewLoginService membuat instance service baru
func NewLoginService(repo Repository) Service {
    return &service{repo}
}

func (s *service) Authenticate(req LoginRequest) (*LoginResponse, error) {
    user, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        // Jangan expose detail error database ke luar
        return nil, errors.New("invalid email or password")
    }

    if !util.CheckPasswordHash(req.Password, user.Password) {
        return nil, errors.New("invalid email or password")
    }

    token, err := util.GenerateToken(user.ID)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    return &LoginResponse{Token: token}, nil
}
```

### Yang Boleh dan Tidak Boleh

| Boleh | Tidak Boleh |
|---|---|
| Panggil method repository | Akses `*gorm.DB` langsung |
| Hash / check password | Import `fiber` |
| Generate token | Return error mentah dari DB |
| Validasi logic bisnis | Akses `fiber.Ctx` |
| Return pesan error ramah pengguna | Import controller/router |

---

## Layer Controller

Controller bertanggung jawab untuk **menangani HTTP request dan response**. Menggunakan service melalui interface.

### Aturan

- Nama file: `{feature}_controller.go`
- Struct `Controller` ditulis dengan huruf besar (exported)
- Struct `Controller` hanya boleh menyimpan `Service` (interface)
- Selalu validasi request body sebelum diteruskan ke service
- HTTP status code harus konsisten (lihat tabel di bawah)
- Response selalu dibungkus dengan key `data` untuk sukses, `error` untuk gagal

### HTTP Status Code Standar

| Status | Keterangan |
|---|---|
| `200` | Sukses |
| `400` | Request tidak valid / validasi gagal |
| `401` | Tidak terautentikasi |
| `403` | Tidak punya akses |
| `404` | Data tidak ditemukan |
| `500` | Error server |

### Struktur Wajib

```go
package login

import (
    "starter-wahcah-be/internal/util"
    "github.com/gofiber/fiber/v2"
)

// Controller menangani HTTP request untuk login
type Controller struct {
    service Service
}

// NewLoginController membuat instance controller baru
func NewLoginController(service Service) *Controller {
    return &Controller{service}
}

// Login menangani POST /auth/login
func (c *Controller) Login(ctx *fiber.Ctx) error {
    // 1. Parse request body
    var req LoginRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }

    // 2. Validasi struct
    if errs := util.ValidateStruct(req); errs != nil {
        return ctx.Status(400).JSON(fiber.Map{"validation": errs})
    }

    // 3. Panggil service
    res, err := c.service.Authenticate(req)
    if err != nil {
        return ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
    }

    // 4. Return response
    return ctx.JSON(fiber.Map{"data": res})
}
```

### Format Response Standar

```json
// Sukses
{ "data": { ... } }

// Validasi gagal
{ "validation": [{ "field": "Email", "tag": "email" }] }

// Error
{ "error": "pesan error" }
```

### Yang Boleh dan Tidak Boleh

| Boleh | Tidak Boleh |
|---|---|
| Parse body, query, params | Akses database langsung |
| Validasi request | Logic bisnis |
| Panggil method service | Import `gorm` |
| Set HTTP status | Hash password |
| Return JSON response | Import repository |

---

## Layer Router

Router bertanggung jawab untuk **mendaftarkan route** dan **menyambungkan dependency** (repository -> service -> controller).

### Aturan

- Nama file: `{feature}_router.go`
- Nama fungsi selalu `InitRoutes`
- Inisialisasi dependency (repo, service, controller) dilakukan di sini
- Grouping route menggunakan nama domain, contoh `/auth`, `/users`
- Middleware spesifik feature didaftarkan di sini, bukan di `router.go` utama

### Struktur Wajib

```go
package login

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

// InitRoutes mendaftarkan semua route untuk login
func InitRoutes(router fiber.Router, db *gorm.DB) {
    repo := NewLoginRepository(db)
    svc  := NewLoginService(repo)
    ctrl := NewLoginController(svc)

    auth := router.Group("/auth")
    auth.Post("/login", ctrl.Login)
}
```

### Mendaftarkan Route Baru ke Router Utama

Setiap `InitRoutes` harus dipanggil di `internal/router/router.go`:

```go
func SetupRoutes(app *fiber.App, db *gorm.DB) {
    api := app.Group("/api")
    v1  := api.Group("/v1")

    login.InitRoutes(v1, db)
    // register.InitRoutes(v1, db)  <- tambahkan module baru di sini
}
```

---

## Unit Test per Layer

### Aturan Umum Test

- Nama file test: `{feature}_{layer}_test.go`
- Package test: `{feature}_test` — black-box testing, bukan `{feature}`
- Nama fungsi test: `Test{Layer}_{Method}_{Skenario}`
- Setiap method wajib punya minimal **3 skenario**: sukses, input tidak valid, data tidak ditemukan
- Mock di-generate otomatis dengan `make mock`, jangan ditulis manual
- Tidak boleh konek ke database sungguhan — selalu gunakan mock

### Penamaan Fungsi Test

```
Test{Layer}_{Method}_{Skenario}

Contoh:
TestRepository_FindByEmail_Success
TestRepository_FindByEmail_NotFound
TestService_Authenticate_Success
TestService_Authenticate_WrongPassword
TestService_Authenticate_EmailNotFound
TestController_Login_Success
TestController_Login_InvalidJSON
TestController_Login_ValidationError
TestController_Login_Unauthorized
TestRouter_POST_AuthLogin_RouteExists
TestRouter_GET_AuthLogin_MethodNotAllowed
```

---

### Test Repository

Repository test memverifikasi bahwa **interface Repository berperilaku sesuai kontrak** menggunakan mock.

```go
package login_test

import (
    "errors"
    "testing"

    "starter-wahcah-be/internal/models"
    "starter-wahcah-be/internal/modules/auth/login/mocks"

    "github.com/stretchr/testify/assert"
)

func TestRepository_FindByEmail_Success(t *testing.T) {
    mockRepo := new(mocks.Repository)

    expected := &models.User{Email: "member@example.com", Password: "hashed"}
    mockRepo.On("FindByEmail", "member@example.com").Return(expected, nil)

    result, err := mockRepo.FindByEmail("member@example.com")

    assert.NoError(t, err)
    assert.Equal(t, expected.Email, result.Email)
    mockRepo.AssertExpectations(t)
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
    mockRepo := new(mocks.Repository)

    mockRepo.On("FindByEmail", "ghost@example.com").Return(nil, errors.New("record not found"))

    result, err := mockRepo.FindByEmail("ghost@example.com")

    assert.Error(t, err)
    assert.Nil(t, result)
    mockRepo.AssertExpectations(t)
}
```

---

### Test Service

Service test memverifikasi **logic bisnis** dengan mock repository.

```go
package login_test

import (
    "errors"
    "testing"

    "starter-wahcah-be/internal/models"
    "starter-wahcah-be/internal/modules/auth/login"
    "starter-wahcah-be/internal/modules/auth/login/mocks"
    "starter-wahcah-be/internal/util"

    "github.com/stretchr/testify/assert"
)

func TestService_Authenticate_Success(t *testing.T) {
    mockRepo := new(mocks.Repository)

    hashed, _ := util.HashPassword("member123")
    mockRepo.On("FindByEmail", "member@example.com").Return(
        &models.User{Email: "member@example.com", Password: hashed}, nil,
    )

    svc := login.NewLoginService(mockRepo)
    res, err := svc.Authenticate(login.LoginRequest{
        Email: "member@example.com", Password: "member123",
    })

    assert.NoError(t, err)
    assert.NotEmpty(t, res.Token)
    mockRepo.AssertExpectations(t)
}

func TestService_Authenticate_EmailNotFound(t *testing.T) {
    mockRepo := new(mocks.Repository)

    mockRepo.On("FindByEmail", "ghost@example.com").Return(nil, errors.New("record not found"))

    svc := login.NewLoginService(mockRepo)
    res, err := svc.Authenticate(login.LoginRequest{
        Email: "ghost@example.com", Password: "pass123",
    })

    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Equal(t, "invalid email or password", err.Error())
    mockRepo.AssertExpectations(t)
}

func TestService_Authenticate_WrongPassword(t *testing.T) {
    mockRepo := new(mocks.Repository)

    hashed, _ := util.HashPassword("correctpassword")
    mockRepo.On("FindByEmail", "member@example.com").Return(
        &models.User{Email: "member@example.com", Password: hashed}, nil,
    )

    svc := login.NewLoginService(mockRepo)
    res, err := svc.Authenticate(login.LoginRequest{
        Email: "member@example.com", Password: "wrongpassword",
    })

    assert.Error(t, err)
    assert.Nil(t, res)
    assert.Equal(t, "invalid email or password", err.Error())
    mockRepo.AssertExpectations(t)
}
```

---

### Test Controller

Controller test memverifikasi **HTTP status dan response body** menggunakan `app.Test()` dari Fiber.

```go
package login_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "starter-wahcah-be/internal/modules/auth/login"
    "starter-wahcah-be/internal/modules/auth/login/mocks"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func setupControllerApp(svc login.Service) *fiber.App {
    app := fiber.New()
    ctrl := login.NewLoginController(svc)
    app.Post("/auth/login", ctrl.Login)
    return app
}

func TestController_Login_Success(t *testing.T) {
    mockSvc := new(mocks.Service)
    mockSvc.On("Authenticate", login.LoginRequest{
        Email: "member@example.com", Password: "member123",
    }).Return(&login.LoginResponse{Token: "mocked-token"}, nil)

    app := setupControllerApp(mockSvc)

    body, _ := json.Marshal(fiber.Map{"email": "member@example.com", "password": "member123"})
    req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(res.Body).Decode(&response)
    data := response["data"].(map[string]interface{})
    assert.Equal(t, "mocked-token", data["token"])
    mockSvc.AssertExpectations(t)
}

func TestController_Login_ValidationError(t *testing.T) {
    mockSvc := new(mocks.Service)
    app := setupControllerApp(mockSvc)

    body, _ := json.Marshal(fiber.Map{"email": "bukan-email", "password": "123"})
    req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    res, _ := app.Test(req)
    assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Login_Unauthorized(t *testing.T) {
    mockSvc := new(mocks.Service)
    mockSvc.On("Authenticate", login.LoginRequest{
        Email: "member@example.com", Password: "wrongpassword",
    }).Return(nil, assert.AnError)

    app := setupControllerApp(mockSvc)

    body, _ := json.Marshal(fiber.Map{"email": "member@example.com", "password": "wrongpassword"})
    req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    res, _ := app.Test(req)
    assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
    mockSvc.AssertExpectations(t)
}
```

---

### Test Router

Router test memverifikasi bahwa **route terdaftar dengan benar** dan method yang salah ditolak.

```go
package login_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "starter-wahcah-be/internal/modules/auth/login"
    "starter-wahcah-be/internal/modules/auth/login/mocks"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func setupRouterApp(svc login.Service) *fiber.App {
    app := fiber.New()
    ctrl := login.NewLoginController(svc)
    auth := app.Group("/auth")
    auth.Post("/login", ctrl.Login)
    return app
}

func TestRouter_POST_AuthLogin_RouteExists(t *testing.T) {
    mockSvc := new(mocks.Service)
    mockSvc.On("Authenticate", login.LoginRequest{
        Email: "member@example.com", Password: "member123",
    }).Return(&login.LoginResponse{Token: "mocked-token"}, nil)

    app := setupRouterApp(mockSvc)

    body, _ := json.Marshal(fiber.Map{"email": "member@example.com", "password": "member123"})
    req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    res, _ := app.Test(req)
    assert.NotEqual(t, http.StatusNotFound, res.StatusCode)
    assert.Equal(t, http.StatusOK, res.StatusCode)
    mockSvc.AssertExpectations(t)
}

func TestRouter_GET_AuthLogin_MethodNotAllowed(t *testing.T) {
    mockSvc := new(mocks.Service)
    app := setupRouterApp(mockSvc)

    req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
    res, _ := app.Test(req)
    assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
```

---

## Test Util

Util test memverifikasi **fungsi-fungsi helper** yang dipakai di seluruh layer.

### password_util_test.go

```go
package util_test

import (
    "testing"

    "starter-wahcah-be/internal/util"

    "github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
    hash, err := util.HashPassword("secret123")

    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
    assert.NotEqual(t, "secret123", hash)
}

func TestHashPassword_ProducesDifferentHashEachTime(t *testing.T) {
    hash1, _ := util.HashPassword("secret123")
    hash2, _ := util.HashPassword("secret123")

    assert.NotEqual(t, hash1, hash2)
}

func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
    hash, _ := util.HashPassword("secret123")

    result := util.CheckPasswordHash("secret123", hash)

    assert.True(t, result)
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
    hash, _ := util.HashPassword("secret123")

    result := util.CheckPasswordHash("wrongpassword", hash)

    assert.False(t, result)
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
    hash, _ := util.HashPassword("secret123")

    result := util.CheckPasswordHash("", hash)

    assert.False(t, result)
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
    result := util.CheckPasswordHash("secret123", "")

    assert.False(t, result)
}
```

### jwt_util_test.go

```go
package util_test

import (
    "os"
    "testing"
    "time"

    "starter-wahcah-be/internal/util"

    "github.com/golang-jwt/jwt/v5"
    "github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")

    token, err := util.GenerateToken(1)

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}

func TestGenerateToken_DifferentUserID(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")

    token1, _ := util.GenerateToken(1)
    token2, _ := util.GenerateToken(2)

    assert.NotEqual(t, token1, token2)
}

func TestGenerateToken_ContainsCorrectUserID(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")

    token, err := util.GenerateToken(99)
    assert.NoError(t, err)

    parsed, parseErr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    assert.NoError(t, parseErr)
    assert.True(t, parsed.Valid)

    claims := parsed.Claims.(jwt.MapClaims)
    assert.Equal(t, float64(99), claims["user_id"])
}

func TestGenerateToken_ContainsExpiry(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")

    token, _ := util.GenerateToken(1)

    parsed, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    claims := parsed.Claims.(jwt.MapClaims)
    exp := int64(claims["exp"].(float64))
    assert.Greater(t, exp, time.Now().Unix())
}

func TestGenerateToken_InvalidWithWrongSecret(t *testing.T) {
    os.Setenv("JWT_SECRET", "correct-secret")

    token, _ := util.GenerateToken(1)

    _, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return []byte("wrong-secret"), nil
    })

    assert.Error(t, err)
}
```

### validation_util_test.go

```go
package util_test

import (
    "testing"

    "starter-wahcah-be/internal/util"

    "github.com/stretchr/testify/assert"
)

type testLoginPayload struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=6"`
}

func TestValidateStruct_AllValid(t *testing.T) {
    payload := testLoginPayload{
        Email:    "user@example.com",
        Password: "secret123",
    }

    errs := util.ValidateStruct(payload)

    assert.Nil(t, errs)
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
    payload := testLoginPayload{
        Email:    "bukan-email",
        Password: "secret123",
    }

    errs := util.ValidateStruct(payload)

    assert.NotNil(t, errs)
    assert.Equal(t, 1, len(errs))
    assert.Equal(t, "Email", errs[0].Field)
    assert.Equal(t, "email", errs[0].Tag)
}

func TestValidateStruct_PasswordTooShort(t *testing.T) {
    payload := testLoginPayload{
        Email:    "user@example.com",
        Password: "123",
    }

    errs := util.ValidateStruct(payload)

    assert.NotNil(t, errs)
    assert.Equal(t, 1, len(errs))
    assert.Equal(t, "Password", errs[0].Field)
    assert.Equal(t, "min", errs[0].Tag)
}

func TestValidateStruct_MultipleErrors(t *testing.T) {
    payload := testLoginPayload{
        Email:    "bukan-email",
        Password: "123",
    }

    errs := util.ValidateStruct(payload)

    assert.NotNil(t, errs)
    assert.Equal(t, 2, len(errs))
}

func TestValidateStruct_EmptyFields(t *testing.T) {
    payload := testLoginPayload{
        Email:    "",
        Password: "",
    }

    errs := util.ValidateStruct(payload)

    assert.NotNil(t, errs)
    assert.Equal(t, 2, len(errs))

    fields := []string{errs[0].Field, errs[1].Field}
    assert.Contains(t, fields, "Email")
    assert.Contains(t, fields, "Password")
}
```

---

## Test Middleware

Middleware test memverifikasi **behavior Protected dan OptionalAuth** tanpa konek ke database.

### auth_middleware_test.go

```go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "starter-wahcah-be/internal/middleware"
    "starter-wahcah-be/internal/util"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func generateTestToken(userID uint) string {
    os.Setenv("JWT_SECRET", "test-secret-key")
    token, _ := util.GenerateToken(userID)
    return token
}

func setupProtectedApp() *fiber.App {
    app := fiber.New()
    app.Get("/protected", middleware.Protected(), func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"user_id": c.Locals("user_id")})
    })
    return app
}

func setupOptionalApp() *fiber.App {
    app := fiber.New()
    app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"user_id": c.Locals("user_id")})
    })
    return app
}

// Protected: token valid
func TestProtected_WithValidToken_ShouldPass(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupProtectedApp()

    token := generateTestToken(1)
    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)
}

// Protected: tanpa token
func TestProtected_WithoutToken_ShouldReturn401(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupProtectedApp()

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// Protected: token tidak valid
func TestProtected_WithInvalidToken_ShouldReturn401(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupProtectedApp()

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer token-tidak-valid")

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// Protected: secret salah
func TestProtected_WithWrongSecret_ShouldReturn401(t *testing.T) {
    os.Setenv("JWT_SECRET", "secret-A")
    token := generateTestToken(1)

    os.Setenv("JWT_SECRET", "secret-B")
    app := setupProtectedApp()

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// OptionalAuth: token valid
func TestOptionalAuth_WithValidToken_ShouldPass(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupOptionalApp()

    token := generateTestToken(7)
    req := httptest.NewRequest(http.MethodGet, "/optional", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)
}

// OptionalAuth: tanpa token tetap lanjut
func TestOptionalAuth_WithoutToken_ShouldStillPass(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupOptionalApp()

    req := httptest.NewRequest(http.MethodGet, "/optional", nil)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)
}

// OptionalAuth: token invalid tetap lanjut
func TestOptionalAuth_WithInvalidToken_ShouldStillPass(t *testing.T) {
    os.Setenv("JWT_SECRET", "test-secret-key")
    app := setupOptionalApp()

    req := httptest.NewRequest(http.MethodGet, "/optional", nil)
    req.Header.Set("Authorization", "Bearer token-tidak-valid")

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)
}

// OptionalAuth: secret salah tetap lanjut
func TestOptionalAuth_WithWrongSecret_ShouldStillPass(t *testing.T) {
    os.Setenv("JWT_SECRET", "secret-A")
    token := generateTestToken(1)

    os.Setenv("JWT_SECRET", "secret-B")
    app := setupOptionalApp()

    req := httptest.NewRequest(http.MethodGet, "/optional", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    res, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, res.StatusCode)
}
```

---

## Integration Test Repository

Integration test menggunakan **SQLite in-memory** — tidak konek ke MySQL sungguhan, aman dijalankan di CI/CD.

### Dependency Tambahan

```bash
go get gorm.io/driver/sqlite
go mod tidy
```

### login_repository_integration_test.go

```go
package login_test

import (
    "testing"

    "starter-wahcah-be/internal/models"
    "starter-wahcah-be/internal/modules/auth/login"

    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }

    if err := db.AutoMigrate(&models.User{}); err != nil {
        t.Fatalf("Failed to migrate test database: %v", err)
    }

    return db
}

func seedTestUser(db *gorm.DB, email, password string) *models.User {
    user := &models.User{Email: email, Password: password}
    db.Create(user)
    return user
}

func TestIntegration_Repository_FindByEmail_Success(t *testing.T) {
    db := setupTestDB(t)
    seedTestUser(db, "member@example.com", "hashedpassword")

    repo := login.NewLoginRepository(db)
    result, err := repo.FindByEmail("member@example.com")

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "member@example.com", result.Email)
}

func TestIntegration_Repository_FindByEmail_NotFound(t *testing.T) {
    db := setupTestDB(t)

    repo := login.NewLoginRepository(db)
    _, err := repo.FindByEmail("ghost@example.com")

    assert.Error(t, err)
}

func TestIntegration_Repository_FindByEmail_ReturnsCorrectUser(t *testing.T) {
    db := setupTestDB(t)
    seedTestUser(db, "user1@example.com", "hash1")
    seedTestUser(db, "user2@example.com", "hash2")

    repo := login.NewLoginRepository(db)
    result, err := repo.FindByEmail("user2@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "user2@example.com", result.Email)
    assert.Equal(t, "hash2", result.Password)
}

func TestIntegration_Repository_FindByEmail_ReturnPasswordHash(t *testing.T) {
    db := setupTestDB(t)
    seedTestUser(db, "member@example.com", "$2a$14$hashedpasswordexample")

    repo := login.NewLoginRepository(db)
    result, err := repo.FindByEmail("member@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "$2a$14$hashedpasswordexample", result.Password)
}
```

---

## Checklist Membuat Module Baru

Ikuti urutan ini setiap kali membuat feature baru:

```
[ ] 1.  Buat folder internal/modules/{domain}/{feature}/
[ ] 2.  Buat {feature}_dto.go          - definisikan Request & Response struct
[ ] 3.  Buat {feature}_repository.go   - interface + implementasi + constructor
[ ] 4.  Buat {feature}_service.go      - interface + implementasi + constructor
[ ] 5.  Buat {feature}_controller.go   - struct Controller + method handler
[ ] 6.  Buat {feature}_router.go       - fungsi InitRoutes + wiring dependency
[ ] 7.  Daftarkan InitRoutes di internal/router/router.go
[ ] 8.  Jalankan make mock              - generate mock otomatis
[ ] 9.  Buat {feature}_repository_test.go
[ ] 10. Buat {feature}_service_test.go
[ ] 11. Buat {feature}_controller_test.go
[ ] 12. Buat {feature}_router_test.go
[ ] 13. Buat {feature}_repository_integration_test.go
[ ] 14. Jalankan make test              - pastikan semua PASS
```

---

## Perintah Penting

```bash
make run                                              # Jalankan server
make dev                                              # Live reload dengan air
make mock                                             # Generate semua mock di internal/modules
make test                                             # Semua test (unit + integration + util + middleware)
make test-unit                                        # Hanya test modules
make test-util                                        # Hanya test util dan middleware
make test-integration                                 # Hanya integration test
make test-pkg PKG=./internal/modules/auth/login/...  # Test satu package spesifik
make test-coverage                                    # Test + laporan coverage HTML
make fmt                                              # Format semua kode
make tidy                                             # Tidy dependencies
make help                                             # Lihat semua perintah
```