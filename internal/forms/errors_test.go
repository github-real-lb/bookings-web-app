package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_Add(t *testing.T) {
	form := createRandomForm(t)
	form.Errors.Add("key1", "Error")
	assert.Len(t, form.Errors, 1)
	assert.Equal(t, "Error", form.Errors.Get("key1"))
}

func TestErrors_Get(t *testing.T) {
	form := createRandomForm(t)
	form.Errors.Add("key1", "Error1")
	form.Errors.Add("key1", "Error2")

	assert.Equal(t, "Error1", form.Errors.Get("key1"))
	assert.Empty(t, form.Errors.Get("key2"))
}
