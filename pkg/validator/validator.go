package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator оборачивает validator для использования с Gin
type CustomValidator struct {
	validator *validator.Validate
}

// New создает новый экземпляр валидатора
func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate выполняет валидацию структуры
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return cv.formatValidationError(err)
	}
	return nil
}

// formatValidationError форматирует ошибки валидации в человекочитаемый формат
func (cv *CustomValidator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMessages []string
		for _, e := range validationErrors {
			errMessages = append(errMessages, cv.formatFieldError(e))
		}
		return fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}
	return err
}

// formatFieldError форматирует ошибку конкретного поля
func (cv *CustomValidator) formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("поле '%s' обязательно для заполнения", field)
	case "min":
		return fmt.Sprintf("поле '%s' должно содержать минимум %s символов", field, e.Param())
	case "max":
		return fmt.Sprintf("поле '%s' должно содержать максимум %s символов", field, e.Param())
	case "oneof":
		return fmt.Sprintf("поле '%s' должно быть одним из: %s", field, e.Param())
	default:
		return fmt.Sprintf("поле '%s' не прошло валидацию по правилу '%s'", field, e.Tag())
	}
}
