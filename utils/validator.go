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
	"unique":   "Nilai ini sudah digunakan, mohon gunakan nilai yang berbeda.",
	"len":      "Kolom ini harus berisi %s karakter.",
	"gte":      "Nilai harus lebih besar atau sama dengan %s.",
	"lte":      "Nilai harus lebih kecil atau sama dengan %s.",
	"rfc3339":  "Format tanggal tidak valid.",
	"datetime": "Format tanggal tidak valid.",
	"numeric":  "Nilai harus berupa angka.",
	"uuid":     "Format UUID tidak valid.",
	"email":    "Format email tidak valid.",
	"oneof":    "Nilai hanya bisa salah satu dari: %s.",
	"min":      "Nilai minimal adalah %s.",
	"max":      "Nilai maksimal adalah %s.",

	//// Users Table
	//"users.username.required":      "Username wajib diisi.",
	//"users.username.unique":        "Username sudah digunakan, mohon pilih username lain.",
	//"users.username.len":           "Username harus memiliki panjang maksimal 100 karakter.",
	//"users.password_hash.required": "Password wajib diisi.",
	//"users.is_active.boolean":      "Status harus berupa nilai true atau false.",
	//
	//// Roles Table
	//"roles.name.required": "Nama peran wajib diisi.",
	//"roles.name.unique":   "Nama peran sudah ada.",
	//"roles.name.len":      "Nama peran harus memiliki panjang maksimal 50 karakter.",
	//
	//// Employees Table
	//"employees.name.required":       "Nama karyawan wajib diisi.",
	//"employees.name.len":            "Nama karyawan tidak boleh lebih dari 150 karakter.",
	//"employees.position.required":   "Posisi karyawan wajib diisi.",
	//"employees.position.len":        "Posisi tidak boleh lebih dari 100 karakter.",
	//"employees.phone.required":      "Nomor telepon wajib diisi.",
	//"employees.phone.len":           "Nomor telepon harus terdiri dari maksimal 20 karakter.",
	//"employees.nik.required":        "NIK wajib diisi.",
	//"employees.nik.unique":          "NIK sudah digunakan.",
	//"employees.nik.len":             "NIK harus terdiri dari maksimal 20 karakter.",
	//"employees.hired_date.required": "Tanggal diterima bekerja wajib diisi.",
	//"employees.hired_date.datetime": "Format tanggal salah",
	//
	//// Wages Table
	//"wages.total_amount.required": "Total gaji wajib diisi.",
	//"wages.total_amount.numeric":  "Total gaji harus berupa angka.",
	//"wages.period_start.required": "Periode mulai wajib diisi.",
	//"wages.period_start.datetime": "Format tanggal salah",
	//"wages.period_end.required":   "Periode selesai wajib diisi.",
	//"wages.period_end.datetime":   "Format tanggal salah",
	//
	//// Financial Transactions Table
	//"financial_transactions.amount.required":           "Jumlah transaksi wajib diisi.",
	//"financial_transactions.amount.numeric":            "Jumlah transaksi harus berupa angka.",
	//"financial_transactions.type.required":             "Tipe transaksi wajib diisi.",
	//"financial_transactions.type.oneof":                "Tipe transaksi hanya bisa salah satu dari: kredit, debit.",
	//"financial_transactions.transaction_date.required": "Tanggal transaksi wajib diisi.",
	//"financial_transactions.transaction_date.datetime": "Format tanggal tidak valid..",
	//
	//// Inventory Table
	//"inventory.name.required":             "Nama inventaris wajib diisi.",
	//"inventory.name.len":                  "Nama inventaris tidak boleh lebih dari 255 karakter.",
	//"inventory.quantity.required":         "Jumlah stok wajib diisi.",
	//"inventory.quantity.numeric":          "Jumlah stok harus berupa angka.",
	//"inventory.minimum_quantity.required": "Stok minimum wajib diisi.",
	//"inventory.minimum_quantity.numeric":  "Stok minimum harus berupa angka.",
	//
	//// Orders Table
	//"orders.customer_name.required": "Nama pelanggan wajib diisi.",
	//"orders.customer_name.len":      "Nama pelanggan tidak boleh lebih dari 255 karakter.",
	//"orders.status.required":        "Status pesanan wajib diisi.",
	//"orders.status.oneof":           "Status pesanan hanya bisa salah satu dari: pending, selesai, batal.",
	//"orders.total_amount.required":  "Total jumlah pesanan wajib diisi.",
	//"orders.total_amount.numeric":   "Total jumlah pesanan harus berupa angka.",
	//
	//// Order Items Table
	//"order_items.quantity.required": "Jumlah item wajib diisi.",
	//"order_items.quantity.numeric":  "Jumlah item harus berupa angka.",
	//"order_items.price.required":    "Harga wajib diisi.",
	//"order_items.price.numeric":     "Harga harus berupa angka.",
	//
	//// Inventory Logs Table
	//"inventory_logs.action.required":   "Aksi wajib diisi.",
	//"inventory_logs.action.oneof":      "Aksi hanya bisa salah satu dari: add, remove.",
	//"inventory_logs.quantity.required": "Jumlah wajib diisi.",
	//"inventory_logs.quantity.numeric":  "Jumlah harus berupa angka.",
}

// Validator instance
var validate = validator.New()

// RegisterCustomValidators registers custom validation tags
func RegisterCustomValidators() {
	validate.RegisterValidation("rfc3339", func(fl validator.FieldLevel) bool {
		// Parse the string as RFC3339
		_, err := time.Parse(time.RFC3339, fl.Field().String())
		return err == nil
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

		val := reflect.ValueOf(s)
		for _, fieldErr := range validationErrors {
			field, _ := val.Type().Elem().FieldByName(fieldErr.StructField())
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = fieldErr.Field()
			} else {
				jsonTag = strings.Split(jsonTag, ",")[0]
			}

			fieldName := fieldErr.StructField()
			tag := fieldErr.Tag()
			key := fmt.Sprintf("%s.%s", fieldName, tag)

			message := validationMessages[key]
			if message == "" {
				message = validationMessages[tag]
			}

			if fieldErr.Param() != "" {
				message = fmt.Sprintf(message, fieldErr.Param())
			}

			errorsVal[jsonTag] = message
		}

		return errorsVal
	}

	return nil
}

// ValidateRFC3339Date check validation for date time format
func ValidateRFC3339Date(dateStr string) bool {
	_, err := time.Parse(time.RFC3339, dateStr)
	return err == nil
}
