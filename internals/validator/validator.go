package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// check if Validator is valid (FieldErrors is empty)
func (v *Validator) Valid() bool {
	return len(v.FieldErrors)==0
}

// adds an error message to FieldError with key
func (v *Validator) AddFieldError(key, message string) {
	// initialize empty map if FieldError is nil
	if v.FieldErrors==nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// add error to FieldError only if validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// helper functions for validations
// check if string is empty
func NotBlank(value string) bool {
	return strings.TrimSpace(value)!=""
}

// check if string len is less than limit
func MaxLen(value string, limit int) bool {
	return utf8.RuneCountInString(value)<=limit
}

// check if value is one of permitted values
func PermittedInt(value int, permittedValues ...int) bool {
	for _, pv := range permittedValues {
		if value==pv {
			return true
		}
	}
	return false
}
