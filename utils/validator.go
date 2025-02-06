package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var validationMessages = map[string]string{
	// General
	"required": "Kolom ini wajib diisi.",
	"len":      "Kolom ini harus berisi %s karakter.",
	"gte":      "Harus lebih besar atau sama dengan %s.",
	"lte":      "Harus lebih kecil atau sama dengan %s.",
	"rfc3339":  "Format tanggal tidak valid.",
	"datetime": "Format tanggal tidak valid.",
	"numeric":  "Harus berupa angka.",
	"uuid":     "Format UUID tidak valid.",
	"email":    "Format email tidak valid.",
	"oneof":    "Hanya bisa salah satu dari: %s.",
	"min":      "Minimal isi adalah %s.",
	"max":      "Maksimal isi adalah %s.",

	"gt":       "Harus lebih besar dari %s.",
	"lt":       "Harus lebih kecil dari %s.",
	"eqfield":  "Harus sama dengan %s.",
	"nefield":  "Nilai tidak boleh sama dengan %s.",
	"gtfield":  "Harus lebih besar dari %s.",
	"ltfield":  "Harus lebih kecil dari %s.",
	"gtefield": "Harus lebih besar atau sama dengan %s.",
	"ltefield": "Harus lebih kecil atau sama dengan %s.",
}

// Validator instance
var validate = validator.New()

// RegisterCustomValidators registers custom validation tags
func RegisterCustomValidators() {
	// RFC3339 datetime format
	validate.RegisterValidation("rfc3339", func(fl validator.FieldLevel) bool {
		_, err := time.Parse(time.RFC3339, fl.Field().String())
		return err == nil
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// DecodeAndValidate decodes JSON from the request body and validates the struct
func DecodeAndValidate(r *http.Request, v interface{}) map[string]string {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return map[string]string{"general": "Data tidak valid"}
	}

	validationErrors := ValidateStruct(v)
	return validationErrors
}

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorsVal := make(map[string]string)

		for _, fieldErr := range validationErrors {
			tag := fieldErr.Tag()
			key := fmt.Sprintf("%s.%s", fieldErr.Field(), tag)

			message := validationMessages[key]
			if message == "" {
				message = validationMessages[tag]
			}

			if fieldErr.Param() != "" {
				message = fmt.Sprintf(message, fieldErr.Param())
			}

			errorsVal[fieldErr.Field()] = message
		}

		return errorsVal
	}

	return nil
}
