package util_test

import (
	"testing"

	"starter-wahcah-be/internal/util"

	"github.com/stretchr/testify/assert"
)

// =========================================
// Struct dummy untuk keperluan test
// =========================================

type testLoginPayload struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

type testRequiredOnly struct {
	Name string `validate:"required"`
	Age  int    `validate:"required,min=1"`
}

// =========================================
// ValidateStruct
// =========================================

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
	// Harus ada 2 error: email tidak valid + password terlalu pendek
	assert.Equal(t, 2, len(errs))
}

func TestValidateStruct_EmptyFields(t *testing.T) {
	payload := testLoginPayload{
		Email:    "",
		Password: "",
	}

	errs := util.ValidateStruct(payload)

	assert.NotNil(t, errs)
	// Email required + Password required
	assert.Equal(t, 2, len(errs))

	fields := []string{errs[0].Field, errs[1].Field}
	assert.Contains(t, fields, "Email")
	assert.Contains(t, fields, "Password")
}

func TestValidateStruct_RequiredFieldMissing(t *testing.T) {
	payload := testRequiredOnly{
		Name: "",
		Age:  0,
	}

	errs := util.ValidateStruct(payload)

	assert.NotNil(t, errs)
	assert.GreaterOrEqual(t, len(errs), 1)
}

func TestValidateStruct_ErrorResponseFields(t *testing.T) {
	payload := testLoginPayload{
		Email:    "invalid",
		Password: "secret123",
	}

	errs := util.ValidateStruct(payload)

	assert.NotNil(t, errs)
	// Pastikan struct ErrorResponse memiliki Field dan Tag yang terisi
	assert.NotEmpty(t, errs[0].Field)
	assert.NotEmpty(t, errs[0].Tag)
}
