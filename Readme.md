# Build dockerfile

```bash
podman build -t my-go-starter-img -f .devcontainer/Dockerfile .
```

# 1. Inisialisasi Module (jika belum)

```bash
go mod init starter-wahcah-be
```

# 2. Download Library Utama

```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/mysql             # Driver MySQL (Sesuai container Anda)
go get golang.org/x/crypto/bcrypt       # Untuk Hash Password
go get github.com/golang-jwt/jwt/v5     # Untuk Token
go get github.com/joho/godotenv         # Untuk .env
go get github.com/go-playground/validator/v10 # Untuk Validasi Input
```

# 3. Rapikan

```bash
go mod tidy
```

# Pindah ke air terminal dahulu jalankan perintah dibawah ini

```bash
podman run -it --rm --name backend-go --network devcontainer_starter-network -v ${PWD}:/app -w /app -p 8080:8080 my-go-starter-img sh
```

<!-- C -->
# jalankan

```bash
air
```

---

# Test endpoint

```url
POST http://localhost:9090/api/auth/register-test
POST http://localhost:9090/api/auth/login
```

request raw body :

```json
{ "email": "admin@example.com", "password": "password123" }
```
