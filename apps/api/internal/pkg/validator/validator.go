package validator

import (
	"fmt"
	"reflect"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *govalidator.Validate
}

func New() *CustomValidator {
	v := govalidator.New()

	// Use json tag name instead of struct field name if available
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return strings.ToLower(fld.Name)
		}
		return name
	})

	return &CustomValidator{validate: v}
}

func (cv *CustomValidator) Validate(s interface{}) map[string]interface{} {
	err := cv.validate.Struct(s)
	if err == nil {
		return nil
	}

	errs, ok := err.(govalidator.ValidationErrors)
	if !ok {
		return map[string]interface{}{
			"general": "data input tidak valid",
		}
	}

	messages := make(map[string]interface{})
	for _, e := range errs {
		field := e.Field()
		messages[field] = translateError(e)
	}

	return messages
}

func translateError(e govalidator.FieldError) string {
	param := e.Param()
	switch e.Tag() {
	case "required":
		return "kolom ini wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("minimal %s karakter", param)
		}
		return fmt.Sprintf("nilai minimal adalah %s", param)
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("maksimal %s karakter", param)
		}
		return fmt.Sprintf("nilai maksimal adalah %s", param)
	case "numeric":
		return "harus berupa angka"
	case "alphanum":
		return "hanya boleh berisi kombinasi huruf dan angka"
	case "len":
		return fmt.Sprintf("panjang karakter harus tepat %s", param)
	case "oneof":
		return fmt.Sprintf("harus merupakan salah satu dari pilihan berikut: %s", param)
	case "gt":
		return fmt.Sprintf("harus lebih besar dari %s", param)
	case "gte":
		return fmt.Sprintf("harus lebih besar atau sama dengan %s", param)
	case "lt":
		return fmt.Sprintf("harus lebih kecil dari %s", param)
	case "lte":
		return fmt.Sprintf("harus lebih kecil atau sama dengan %s", param)
	default:
		return fmt.Sprintf("tidak memenuhi validasi '%s'", e.Tag())
	}
}
