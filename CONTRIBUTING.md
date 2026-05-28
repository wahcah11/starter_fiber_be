# Panduan Penulisan Kode — starter-wahcah-be

Dokumen ini menjelaskan aturan penulisan kode untuk setiap layer, mulai dari penamaan file, interface, struct, hingga unit test. Wajib dibaca sebelum menambahkan module baru.

---

## Daftar Isi

1. [Struktur Folder](#struktur-folder)
2. [Aturan Umum](#aturan-umum)
3. [Layer DTO](#layer-dto)
4. [Layer Repository](#layer-repository)
5. [Layer Service](#layer-service)
6. [Layer Controller](#layer-controller)
7. [Layer Router](#layer-router)
8. [Unit Test](#unit-test)
9. [Checklist Membuat Module Baru](#checklist-membuat-module-baru)

---

## Struktur Folder

Setiap module berada di dalam `internal/modules/{domain}/{feature}/`.

```
internal/
└── modules/
    └── {domain}/               # contoh: auth, user, product
        └── {feature}/          # contoh: login, register, profile
            ├── {feature}_dto.go
            ├── {feature}_repository.go
            ├── {feature}_service.go
            ├── {feature}_controller.go
            ├── {feature}_router.go
            ├── {feature}_repository_test.go
            ├── {feature}_service_test.go
            ├── {feature}_controller_test.go
            ├── {feature}_router_test.go
            └── mocks/
                ├── Repository.go   ← auto-generated mockery
                └── Service.go      ← auto-generated mockery
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

### Penamaan Package

Semua file dalam satu folder feature menggunakan package yang sama, yaitu nama feature-nya.

```go
// Semua file di folder login/ menggunakan:
package login
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
package {feature}

import (
    "{module}/internal/models"
    "gorm.io/gorm"
)

// Repository mendefinisikan kontrak operasi database untuk {feature}
type Repository interface {
    FindByEmail(email string) (*models.User, error)
    // tambahkan method lain di sini
}

// repository adalah implementasi dari interface Repository
type repository struct {
    db *gorm.DB
}

// New{Feature}Repository membuat instance repository baru
func New{Feature}Repository(db *gorm.DB) Repository {
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
package {feature}

import (
    "errors"
    "{module}/internal/util"
)

// Service mendefinisikan kontrak logic bisnis untuk {feature}
type Service interface {
    Authenticate(req LoginRequest) (*LoginResponse, error)
    // tambahkan method lain di sini
}

// service adalah implementasi dari interface Service
type service struct {
    repo Repository
}

// New{Feature}Service membuat instance service baru
func New{Feature}Service(repo Repository) Service {
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
- HTTP status code harus konsisten:
  - `200` → sukses
  - `400` → request tidak valid / validasi gagal
  - `401` → tidak terautentikasi
  - `403` → tidak punya akses
  - `404` → data tidak ditemukan
  - `500` → error server
- Response selalu dibungkus dengan key `data` untuk sukses, `error` untuk gagal

### Struktur Wajib

```go
package {feature}

import (
    "{module}/internal/util"
    "github.com/gofiber/fiber/v2"
)

// Controller menangani HTTP request untuk {feature}
type Controller struct {
    service Service
}

// New{Feature}Controller membuat instance controller baru
func New{Feature}Controller(service Service) *Controller {
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
{ "validation": { "email": "email tidak valid", "password": "minimal 6 karakter" } }

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

Router bertanggung jawab untuk **mendaftarkan route** dan **menyambungkan dependency** (repository → service → controller).

### Aturan

- Nama file: `{feature}_router.go`
- Nama fungsi selalu `InitRoutes`
- Inisialisasi dependency (repo, service, controller) dilakukan di sini
- Grouping route menggunakan nama domain, contoh `/auth`, `/users`
- Middleware spesifik feature didaftarkan di sini, bukan di `router.go` utama

### Struktur Wajib

```go
package {feature}

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

// InitRoutes mendaftarkan semua route untuk {feature}
func InitRoutes(router fiber.Router, db *gorm.DB) {
    // Inisialisasi dependency
    repo := New{Feature}Repository(db)
    svc  := New{Feature}Service(repo)
    ctrl := New{Feature}Controller(svc)

    // Grouping route
    group := router.Group("/{domain}")
    group.Post("/{endpoint}", ctrl.{Method})
}
```

Contoh:

```go
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
    register.InitRoutes(v1, db)  // ← tambahkan di sini
}
```

---

## Unit Test

### Aturan Umum Test

- Nama file test: `{feature}_{layer}_test.go`
- Package test: `{feature}_test` (bukan `{feature}`) — ini black-box testing
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

// Helper: buat app fiber dengan controller yang di-inject mock service
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
```

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

// Helper: buat app fiber dengan route yang di-inject mock service
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

## Checklist Membuat Module Baru

Ikuti urutan ini setiap kali membuat feature baru:

```
[ ] 1. Buat folder internal/modules/{domain}/{feature}/
[ ] 2. Buat {feature}_dto.go         — definisikan Request & Response struct
[ ] 3. Buat {feature}_repository.go  — interface + implementasi + constructor
[ ] 4. Buat {feature}_service.go     — interface + implementasi + constructor
[ ] 5. Buat {feature}_controller.go  — struct Controller + method handler
[ ] 6. Buat {feature}_router.go      — fungsi InitRoutes + wiring dependency
[ ] 7. Daftarkan InitRoutes di internal/router/router.go
[ ] 8. Jalankan make mock             — generate mock otomatis
[ ] 9. Buat {feature}_repository_test.go
[ ] 10. Buat {feature}_service_test.go
[ ] 11. Buat {feature}_controller_test.go
[ ] 12. Buat {feature}_router_test.go
[ ] 13. Jalankan make test            — pastikan semua PASS
```

---

## Perintah Penting

```bash
make run            # Jalankan server
make dev            # Jalankan dengan live reload (air)
make mock           # Generate semua mock di internal/modules
make test           # Jalankan semua test di internal/modules
make test-coverage  # Test + laporan coverage HTML
make test-pkg PKG=./internal/modules/auth/login/...  # Test satu package
make fmt            # Format semua kode
make tidy           # Tidy dependencies
make help           # Lihat semua perintah
```