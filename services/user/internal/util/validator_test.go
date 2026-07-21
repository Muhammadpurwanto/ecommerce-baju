package util

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type TestUser struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestFormatValidationErrors(t *testing.T) {
	validate := validator.New()

	t.Run("Valid Struct", func(t *testing.T) {
		user := TestUser{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   20,
		}
		err := validate.Struct(user)
		if err != nil {
			t.Errorf("Diharapkan tidak ada error validasi, namun mendapat: %v", err)
		}
	})

	t.Run("Invalid Struct Fields", func(t *testing.T) {
		user := TestUser{
			Name:  "",              // Error: required
			Email: "invalid-email", // Error: email format
			Age:   15,              // Error: min (di bawah 18)
		}
		err := validate.Struct(user)
		if err == nil {
			t.Fatal("Diharapkan terjadi error validasi, namun tidak ada")
		}

		formattedErrors := FormatValidationErrors(err)

		if len(formattedErrors) != 3 {
			t.Errorf("Diharapkan 3 error validasi, namun mendapat: %d", len(formattedErrors))
		}

		// Pindahkan hasil ke map agar mudah dicek
		errorMap := make(map[string]string)
		for _, fe := range formattedErrors {
			errorMap[fe["field"]] = fe["message"]
		}

		// Cek pesan error field Name
		if msg, ok := errorMap["Name"]; !ok || msg != "Name is required" {
			t.Errorf("Diharapkan 'Name is required', namun mendapat: '%s'", msg)
		}

		// Cek pesan error field Email
		if msg, ok := errorMap["Email"]; !ok || msg != "Email must be a valid email" {
			t.Errorf("Diharapkan 'Email must be a valid email', namun mendapat: '%s'", msg)
		}

		// Cek pesan error field Age (sesuai string bawaan getValidationMessage di validator.go)
		if msg, ok := errorMap["Age"]; !ok || msg != "Age must be at least 18 characters" {
			t.Errorf("Diharapkan 'Age must be at least 18 characters', namun mendapat: '%s'", msg)
		}
	})
}
