package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	form := createRandomForm(t)
	form.Errors.Add("key1", "Error")
	assert.Len(t, form.Errors, 1)
	assert.Equal(t, "Error", form.Errors.Get("key1"))
}

func TestGet(t *testing.T) {
	form := createRandomForm(t)
	form.Errors.Add("key1", "Error1")
	form.Errors.Add("key1", "Error2")

	assert.Equal(t, "Error1", form.Errors.Get("key1"))
	assert.Empty(t, form.Errors.Get("key2"))
}
