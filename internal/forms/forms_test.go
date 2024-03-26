package forms

import (
	"encoding/json"
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

func TestForm_Has(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	for key := range form.Values {
		assert.True(t, form.Has(key))
	}

	// check missing key
	assert.False(t, form.Has(util.RandomName()))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.True(t, form.Has("WhiteSpaces"))
}

func TestForm_GetValue(t *testing.T) {
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

func TestForm_IsEmailValid(t *testing.T) {
	tests := []struct {
		key    string
		value  string
		result bool
	}{
		{"Valid Email1", "john.dow@gmail.com", true},
		{"Valid Email with White Spaces", "   john.dow@gmail.com   ", true},
		{"Invalid Email", "john.dow@", false},
		{"White Spaces", "     ", false},
		{"Empty Field", "", false},
	}

	form := createRandomForm(t)
	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			form.Set(test.key, test.value)
			assert.Equal(t, test.result, form.IsEmailValid(test.key))
		})
	}
}

func TestForm_MarshalJsonFirst(t *testing.T) {
	form := createRandomForm(t)
	jsonData, err := form.MarshalJsonFirst()
	require.NoError(t, err)
	require.NotEmpty(t, jsonData)

	formData := make(map[string]string)
	err = json.Unmarshal(jsonData, &formData)
	require.NoError(t, err)

	for key, value := range formData {
		assert.Equal(t, form.Get(key), value)
	}
}

func TestForm_MarshalJsonAll(t *testing.T) {
	form := createRandomForm(t)
	newJsonData, err := form.MarshalJsonAll()
	require.NoError(t, err)
	require.NotEmpty(t, newJsonData)

	originalJsonData, err := json.Marshal(form.Values)
	require.NoError(t, err)

	assert.Equal(t, originalJsonData, newJsonData)
}

func TestForm_MinLenght(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	assert.True(t, form.MinLenght("key1", 3))

	// check not ok
	assert.False(t, form.MinLenght("key1", 100))

	// check missing key
	assert.False(t, form.MinLenght(util.RandomName(), 3))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.True(t, form.MinLenght("WhiteSpaces", 3))
}

func TestForm_Required(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	assert.True(t, form.Required("key1"))

	// check missing key
	assert.False(t, form.Required(util.RandomName()))

	// check white spaces value
	form = New(url.Values{})
	form.Set("WhiteSpaces", "     ")
	assert.True(t, form.Required("WhiteSpaces"))
}

func TestForm_TrimSpaces(t *testing.T) {
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
	}

	form.TrimSpaces()
	for _, test := range tests {
		assert.Equal(t, test.result, form.Get(test.key))
	}
}

func TestForm_Valid(t *testing.T) {
	form := createRandomForm(t)
	assert.True(t, form.Valid())

	form.Errors.Add("field", "error")
	assert.False(t, form.Valid())
}
