package forms

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// creates a new form with n > 1 values and no errors
// Values keys are names key1, key2, ...
func createRandomForm(t *testing.T) *Form {
	n := 2
	values := url.Values{}

	for i := 1; i <= n; i++ {
		values.Set(fmt.Sprint("key", i), util.RandomName())
	}

	form := New(values)
	require.NotNil(t, form)
	require.NotNil(t, form.Values)
	assert.NotEmpty(t, form.Values)
	assert.Equal(t, n, len(form.Values))
	require.NotNil(t, form.Errors)
	assert.Empty(t, form.Errors)

	return form
}

func TestNew(t *testing.T) {
	createRandomForm(t)
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		key    string
		value  string
		result string
	}{
		{"key1", "  Value  ", "Value"},
		{"key2", "  Value 2  ", "Value 2"},
		{"key3", "    ", ""},
		{"key4", "", ""},
	}

	form := createRandomForm(t)
	for _, test := range tests {
		form.Set(test.key, test.value)
		assert.Equal(t, test.result, form.GetValue(test.key))
	}
}

func TestValid(t *testing.T) {
	form := createRandomForm(t)
	assert.True(t, form.Valid())

	form.Errors.Add("field", "error")
	assert.False(t, form.Valid())
}

func TestHas(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	for key := range form.Values {
		assert.True(t, form.Has(key))
	}

	// check missing key
	assert.False(t, form.Has(util.RandomName()))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.False(t, form.Has("WhiteSpaces"))
}

func TestRequired(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	assert.True(t, form.Required("key1"))

	// check missing key
	assert.False(t, form.Required(util.RandomName()))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.False(t, form.Required("WhiteSpaces"))
}

func TestMinLenght(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	assert.True(t, form.MinLenght("key1", 3))

	// check missing key
	assert.False(t, form.MinLenght(util.RandomName(), 3))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.False(t, form.Required("WhiteSpaces"))
}

func TestIsEmailValid(t *testing.T) {
	tests := []struct {
		key    string
		value  string
		result bool
	}{
		{"Valid Email1", "john.dow@gmail.com", true},
		{"Valid Email2", "   john.dow@gmail.com   ", true},
		{"Invalid Email", "john.dow@", false},
		{"Empty Field", "", false},
	}

	form := createRandomForm(t)
	for _, test := range tests {
		form.Set(test.key, test.value)
		assert.Equal(t, test.result, form.IsEmailValid(test.key))
	}
}
