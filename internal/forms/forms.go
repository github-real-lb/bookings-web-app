package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form is used to hold the data and error of the fields of an html form
type Form struct {
	url.Values
	Errors errors
}

// New initialized a form struct
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
	}
}

// GetValue returns the first value associated with the given field,
// with all leading and trailing white space removed
func (f *Form) GetValue(field string) string {
	return strings.TrimSpace(f.Get(field))
}

// Valid return checks if form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Has checks if the field passed contains data, and returns the result.
// Error message is NOT added to f.Errors.
func (f *Form) Has(field string) bool {
	return !(f.GetValue(field) == "")
}

// Required checks if all fields passed contain data, and returns the result.
// Error messages are added to f.Errors for empty fields.
func (f *Form) Required(fields ...string) bool {
	for _, field := range fields {
		if f.GetValue(field) == "" {
			f.Errors.Add(field, "Required field!")
		}
	}

	return len(f.Errors) == 0
}

// MinLenght checks if the field passed has minimum characters, and returns the result.
// Error message is added to f.Errors.
func (f *Form) MinLenght(field string, lenght int) bool {
	if len(f.GetValue(field)) < lenght {
		f.Errors.Add(field, fmt.Sprintf("Field requires at least %d characters!", lenght))
		return false
	}

	return true
}

// IsEmail checks if the field passed has a valid email, and returns the result.
// Error message is added to f.Errors.
func (f *Form) IsEmailValid(field string) bool {
	if !govalidator.IsEmail(f.GetValue(field)) {
		f.Errors.Add(field, "Invalid email address!")
		return false
	}

	return true
}
