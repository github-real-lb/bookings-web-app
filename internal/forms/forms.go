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

// Has checks if the field passed contains data, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Error message is NOT added to f.Errors in case the field is empty.
func (f *Form) Has(field string) bool {
	return !(f.Get(field) == "")
}

// GetValue returns the first value associated with the given field,
// with all leading and trailing white space removed
func (f *Form) GetValue(field string) string {
	return strings.TrimSpace(f.Get(field))
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

// Marshal returns the data of the first values in each field of f.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
func (f *Form) Marshal() (data map[string]string) {
	// Convert the form data into a map
	data = make(map[string]string)
	for key, values := range f.Values {
		if len(values) != 0 {
			data[key] = values[0]
		}
	}
	return
}

// MinLenght checks if the first value of the field passed has minimum characters, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Error message is added to f.Errors.
func (f *Form) MinLenght(field string, lenght int) bool {
	if len(f.Get(field)) < lenght {
		f.Errors.Add(field, fmt.Sprintf("Field requires at least %d characters!", lenght))
		return false
	}

	return true
}

// Required checks if the first values of all fields passed contain data, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Error messages are added to f.Errors for empty fields.
func (f *Form) Required(fields ...string) bool {
	for _, field := range fields {
		if f.Get(field) == "" {
			f.Errors.Add(field, "Required field!")
		}
	}

	return len(f.Errors) == 0
}

// TrimSpaces removes all leading and trailing white space from the first value of all fields in the form
func (f *Form) TrimSpaces() {
	for key, values := range f.Values {
		if len(values) != 0 {
			f.Values[key][0] = strings.TrimSpace(values[0])
		}
	}
}

// Valid return checks if form is valid
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
