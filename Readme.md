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
go get gorm.io/driver/mysql            
go get golang.org/x/crypto/bcrypt      
go get github.com/golang-jwt/jwt/v5     
go get github.com/joho/godotenv         
go get github.com/go-playground/validator/v10 
```

# 3. Rapikan
```bash
go mod tidy
```

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
{"email": "admin@example.com", "password": "password123"}
```

