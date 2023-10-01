package validator

import (
	"regexp"
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

// check if string is at least limit chars long
func MinLen(value string, limit int) bool {
	return utf8.RuneCountInString(value)>=8	
}

// regex for checking email
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// return true if value matches regular expression
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
