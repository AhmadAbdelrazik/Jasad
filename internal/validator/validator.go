package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// check if form is valid
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// adds an error after validation fails
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists:= v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// validate a field.
// the validation function is passed to ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Chars(value string, min, max int) bool {
	return (utf8.RuneCountInString(value) >= min && utf8.RuneCountInString(value) <= max)
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
