package util_test

import (
	"testing"

	"starter-wahcah-be/internal/util"

	"github.com/stretchr/testify/assert"
)

// =========================================
// HashPassword
// =========================================

func TestHashPassword_Success(t *testing.T) {
	hash, err := util.HashPassword("secret123")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	// Hash tidak boleh sama dengan plain text
	assert.NotEqual(t, "secret123", hash)
}

func TestHashPassword_ProducesDifferentHashEachTime(t *testing.T) {
	// bcrypt menggunakan salt acak, hash yang sama tidak boleh dihasilkan dua kali
	hash1, _ := util.HashPassword("secret123")
	hash2, _ := util.HashPassword("secret123")

	assert.NotEqual(t, hash1, hash2)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	// Password kosong tetap bisa di-hash oleh bcrypt
	hash, err := util.HashPassword("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

// =========================================
// CheckPasswordHash
// =========================================

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

func TestCheckPasswordHash_BothEmpty(t *testing.T) {
	result := util.CheckPasswordHash("", "")

	assert.False(t, result)
}
