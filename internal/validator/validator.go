package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid returns true if FieldErrors doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
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

// AddNonFieldError() helper for adding error messages to the NonFieldErrors[]
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func MinChars(value string, n int) bool { return utf8.RuneCountInString(value) >= n }

// PermittedValue returns true if a value is in a list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Returns a pointer to regexp.Regexp type, or panics.
var EmailRx = regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEqual[T comparable](value1, value2 T) bool {
	if value1 == value2 {
		return true
	}
	return false
}
