package util_test

import (
	"os"
	"testing"
	"time"

	"starter-wahcah-be/internal/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// =========================================
// GenerateToken
// =========================================

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

	// Token untuk user berbeda harus berbeda
	assert.NotEqual(t, token1, token2)
}

func TestGenerateToken_ContainsCorrectUserID(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	token, err := util.GenerateToken(99)
	assert.NoError(t, err)

	// Parse token dan verifikasi user_id di dalam claims
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

	// Pastikan exp ada dan nilainya di masa depan
	exp := int64(claims["exp"].(float64))
	assert.Greater(t, exp, time.Now().Unix())
}

func TestGenerateToken_EmptySecret(t *testing.T) {
	// Jika JWT_SECRET kosong, token tetap dibuat tapi tidak aman
	// Test ini memastikan fungsi tidak panic
	os.Setenv("JWT_SECRET", "")

	token, err := util.GenerateToken(1)

	// Tidak boleh panic, boleh error atau berhasil
	_ = token
	_ = err

	// Reset ke nilai aman
	os.Setenv("JWT_SECRET", "test-secret-key")
}

func TestGenerateToken_InvalidWithWrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "correct-secret")

	token, _ := util.GenerateToken(1)

	// Coba parse dengan secret yang salah
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})

	assert.Error(t, err)
}
