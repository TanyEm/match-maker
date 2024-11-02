package apiserver

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

// validatorTestableFieldLevel is a testable implementation of the validator.FieldLevel interface.
type validatorTestableFieldLevel struct {
	validator.FieldLevel
	Value string
}

func (v *validatorTestableFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(v.Value)
}

func TestISOCountryValidator(t *testing.T) {
	tests := []struct {
		name     string
		country  string
		expected bool
	}{
		{"valid country", "USA", true},
		{"invalid country", "US", false},
		{"empty country", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := &validatorTestableFieldLevel{Value: tt.country}

			if got := ISOCountryValidator(fl); got != tt.expected {
				t.Errorf("ISOCountryValidator() = %v, want %v", got, tt.expected)
			}
		})
	}
}
