package forms

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// creates a new form with n values and no errors
func createRandomForm(t *testing.T) *Form {
	n := 2
	values := url.Values{}

	for i := 1; i <= n; i++ {
		values.Set(fmt.Sprint("key", i), util.RandomName())
	}

	form := New(values)
	require.NotNil(t, form)
	assert.NotEmpty(t, form.Values)
	assert.Equal(t, n, len(form.Values))
	assert.Empty(t, form.Errors)

	return form
}

func TestNew(t *testing.T) {
	createRandomForm(t)
}

func TestValid(t *testing.T) {
	form := createRandomForm(t)
	assert.True(t, form.Valid())

	form.Errors.Add("field", "error")
	assert.False(t, form.Valid())
}
