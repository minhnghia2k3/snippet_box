package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if FieldErrors doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message to the key in FieldErrors map.
func (v *Validator) AddFieldError(key, msg string) {
	// Init
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// Check if exists, then add to the map
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

// CheckField adds an error message to the FieldErrors map only
// if a validation check is not ok!
func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

// NotBlank returns true if a value contains no more than n characters
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
