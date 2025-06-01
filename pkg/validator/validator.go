package validator

import (
	"fmt"
	"regexp"
	"slices"
)

var (
	EmailRX      = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	UsernameRX   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,19}$`)
	NationalIDRX = regexp.MustCompile(`^[0-9]{14}$`)
	MobileRX     = regexp.MustCompile(`^01[0125][0-9]{8}$`)
	YearRX       = regexp.MustCompile(`^(19|20)\d{2}$`)
	URLRX        = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
)

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// In checks if value is inside the given array
func In(value string, array []string) bool {
	return slices.Contains(array, value)
}

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	v.Errors[key] = message
}

func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}

func (v *Validator) CheckStringLength(value string, minimum, maximum int, key string) {
	if len(value) < minimum || len(value) > maximum {
		v.AddError(key, fmt.Sprintf("must be between %d and %d characters long", minimum, maximum))
	}
}
