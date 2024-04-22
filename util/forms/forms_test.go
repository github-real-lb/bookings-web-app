package forms

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/config"
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
	t.Run("OK Values", func(t *testing.T) {
		createRandomForm(t)
	})

	t.Run("OK Nil", func(t *testing.T) {
		f := New(nil)
		require.Empty(t, f.Values)
		require.Empty(t, f.Errors)
	})
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
			assert.Equal(t, test.result, form.CheckEmail(test.key))
		})
	}
}

func TestForm_GetValue(t *testing.T) {
	sString := util.RandomString(12)
	sInt64 := fmt.Sprint(util.RandomInt64(1, 1000))
	sTime := util.RandomDate().Format(config.DateTimeLayout)
	sDate := util.RandomDate().Format(config.DateLayout)

	data := url.Values{}
	data.Add("key1", sString)
	data.Add("key2", sInt64)
	data.Add("key3", sTime)
	data.Add("key4", sDate)

	f := New(data)

	var err error
	var s string
	var i int64
	var dt time.Time

	// test string
	err = f.GetValue("key1", &s)
	require.NoError(t, err)
	assert.Equal(t, sString, s)

	err = f.GetValue("key1", &i)
	require.Error(t, err)

	// test int64
	err = f.GetValue("key2", &i)
	require.NoError(t, err)
	assert.Equal(t, sInt64, fmt.Sprint(i))

	err = f.GetValue("key2", &dt)
	require.Error(t, err)

	// test time
	err = f.GetValue("key3", &dt)
	require.NoError(t, err)
	assert.Equal(t, sTime, dt.Format(config.DateTimeLayout))

	err = f.GetValue("key3", &i)
	require.Error(t, err)

	// test date
	err = f.GetValue("key4", &dt)
	require.NoError(t, err)
	assert.Equal(t, sDate, dt.Format(config.DateLayout))

	err = f.GetValue("key4", &i)
	require.Error(t, err)
}

func TestForm_MinLenght(t *testing.T) {
	form := createRandomForm(t)

	// check ok
	assert.True(t, form.CheckMinLenght("key1", 3))

	// check not ok
	assert.False(t, form.CheckMinLenght("key1", 100))

	// check missing key
	assert.False(t, form.CheckMinLenght(util.RandomName(), 3))

	// check white spaces value
	form.Set("WhiteSpaces", "     ")
	assert.True(t, form.CheckMinLenght("WhiteSpaces", 3))
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
