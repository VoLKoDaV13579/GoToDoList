package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name     string `validate:"required,min=2,max=10"`
	Email    string `validate:"required,email"`
	Age      int    `validate:"min=0,max=120"`
	Category string `validate:"oneof=foo bar baz"`
}

func TestNew(t *testing.T) {
	v := New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validator)
}

func TestCustomValidator_Validate_Success(t *testing.T) {
	v := New()

	t.Run("valid struct", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.NoError(t, err)
	})

	t.Run("valid struct with max length", func(t *testing.T) {
		testData := TestStruct{
			Name:     "1234567890", // Exactly 10 chars
			Email:    "test@test.com",
			Age:      120,
			Category: "bar",
		}

		err := v.Validate(&testData)
		assert.NoError(t, err)
	})

	t.Run("valid struct with min values", func(t *testing.T) {
		testData := TestStruct{
			Name:     "Jo", // Exactly 2 chars
			Email:    "a@b.c",
			Age:      0,
			Category: "baz",
		}

		err := v.Validate(&testData)
		assert.NoError(t, err)
	})
}

func TestCustomValidator_Validate_Errors(t *testing.T) {
	v := New()

	t.Run("missing required field", func(t *testing.T) {
		testData := TestStruct{
			Name:     "", // Missing
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "обязательно")
	})

	t.Run("min length violation", func(t *testing.T) {
		testData := TestStruct{
			Name:     "J", // Too short
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "минимум")
	})

	t.Run("max length violation", func(t *testing.T) {
		testData := TestStruct{
			Name:     "12345678901", // Too long
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "максимум")
	})

	t.Run("invalid email", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "not-an-email", // Invalid email
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
	})

	t.Run("age below minimum", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "john@example.com",
			Age:      -1, // Below minimum
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
	})

	t.Run("age above maximum", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "john@example.com",
			Age:      121, // Above maximum
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
	})

	t.Run("invalid oneof value", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "john@example.com",
			Age:      30,
			Category: "invalid", // Not in allowed values
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "одним из")
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		testData := TestStruct{
			Name:     "", // Missing
			Email:    "", // Missing
			Age:      200,
			Category: "invalid",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		// Check that multiple errors are present
		errStr := err.Error()
		assert.True(t, strings.Contains(errStr, ";") || strings.Count(errStr, "поле") > 1)
	})
}

func TestCustomValidator_FormatFieldError(t *testing.T) {
	v := New()

	t.Run("format required error", func(t *testing.T) {
		testData := TestStruct{
			Name:     "",
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "обязательно для заполнения")
	})

	t.Run("format min error", func(t *testing.T) {
		testData := TestStruct{
			Name:     "J",
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "минимум")
		assert.Contains(t, err.Error(), "2")
	})

	t.Run("format max error", func(t *testing.T) {
		testData := TestStruct{
			Name:     "12345678901",
			Email:    "john@example.com",
			Age:      30,
			Category: "foo",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "максимум")
		assert.Contains(t, err.Error(), "10")
	})

	t.Run("format oneof error", func(t *testing.T) {
		testData := TestStruct{
			Name:     "John",
			Email:    "john@example.com",
			Age:      30,
			Category: "invalid",
		}

		err := v.Validate(&testData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "одним из")
	})
}
