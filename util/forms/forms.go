package forms

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/github-real-lb/bookings-web-app/util/config"
)

// Form is used to hold the data and error of the fields of an html form
type Form struct {
	url.Values `json:"values"`
	Errors     Errors `json:"errors"`
}

// New initialized a form struct
func New(data url.Values) *Form {
	if data == nil {
		data = make(url.Values)
	}

	return &Form{
		Values: data,
		Errors: make(Errors),
	}
}

// CheckDateRange checks if the startDateField is prior to the endDateField, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Any error message is added to f.Errors for the endDateField field.
func (f *Form) CheckDateRange(startDateField string, endDateField string) bool {
	startDate, err := time.Parse(config.DateLayout, f.Get(startDateField))
	if err != nil {
		f.Errors.Add(endDateField, "Invalid date. Please enter date in the following format: YYYY-MM-DD.")
		return false
	}

	endDate, err := time.Parse(config.DateLayout, f.Get(endDateField))
	if err != nil {
		f.Errors.Add(endDateField, "Invalid date. Please enter date in the following format: YYYY-MM-DD.")
		return false
	}

	if endDate.Before(startDate) {
		f.Errors.Add(endDateField, "End date cannot be prior to start date.")
		return false
	}

	return true
}

// CheckEmail checks if the field passed has a valid email, and returns the result.
// Error message is added to f.Errors.
func (f *Form) CheckEmail(field string) bool {
	if !govalidator.IsEmail(strings.TrimSpace(f.Get(field))) {
		f.Errors.Add(field, "Invalid email address!")
		return false
	}

	return true
}

// CheckMinLenght checks if the first value of the field passed has minimum characters, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Error message is added to f.Errors.
func (f *Form) CheckMinLenght(field string, lenght int) bool {
	if len(f.Get(field)) < lenght {
		f.Errors.Add(field, fmt.Sprintf("Field requires at least %d characters!", lenght))
		return false
	}

	return true
}

// GetValue parse the data in field into target
func (f *Form) GetValue(field string, target any) error {
	src := f.Get(field)
	if src == "" {
		return nil
	}

	var err error = nil

	switch p := target.(type) {
	case *string:
		*p = src
	case *int:
		*p, err = strconv.Atoi(src)
	case *int64:
		*p, err = strconv.ParseInt(src, 10, 64)
	case *time.Time:
		if len(src) == 10 {
			*p, err = time.Parse(config.DateLayout, src)
		} else {
			*p, err = time.Parse(config.DateTimeLayout, src)
		}
	default:
		err = fmt.Errorf("unsupported type %T", p)
	}

	return err
}

// Has checks if the field passed contains data, and returns the result.
// Run TrimSpaces before to remove leading and trailing white spaces if needed.
// Error message is NOT added to f.Errors in case the field is empty.
func (f *Form) Has(field string) bool {
	return !(f.Get(field) == "")
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
