package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=owner kasir operator"`
	Quantity int    `json:"quantity" validate:"gte=1,lte=100"`
}

func TestValidator_Valid(t *testing.T) {
	v := New()
	req := sampleRequest{
		Name:     "Budi",
		Email:    "budi@example.com",
		Role:     "kasir",
		Quantity: 5,
	}

	errs := v.Validate(req)
	assert.Nil(t, errs)
}

func TestValidator_Invalid(t *testing.T) {
	v := New()
	req := sampleRequest{
		Name:     "Al",
		Email:    "invalid-email",
		Role:     "admin_palsu",
		Quantity: 0,
	}

	errs := v.Validate(req)
	assert.NotNil(t, errs)
	assert.Equal(t, "minimal 3 karakter", errs["name"])
	assert.Equal(t, "format email tidak valid", errs["email"])
	assert.Contains(t, errs["role"], "harus merupakan salah satu dari pilihan")
	assert.Contains(t, errs["quantity"], "harus lebih besar atau sama dengan 1")
}
