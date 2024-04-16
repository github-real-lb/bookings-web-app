package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerError_Error(t *testing.T) {
	// create an error
	err := ServerError{
		Prompt: "test prompt",
		URL:    "/test_url",
		Err:    errors.New("test error"),
	}

	assert.Equal(t, "\ttest error\n\tPROMPT: test prompt\n\tURL: /test_url\n", err.Error())
}

func TestCreateServerError(t *testing.T) {
	tests := []struct {
		ErrorType ErrorType
		Prompt    string
		URL       string
		Err       error
	}{
		{
			ErrorMissingReservation,
			"No reservation exists. Please make a reservation.",
			"/test",
			errors.New("wrong routing"),
		},
		{
			ErrorParseForm,
			"Unable to parse form.",
			"/test",
			errors.New("This is a test error"),
		},
		{
			ErrorRenderTemplate,
			"Unable to render template.",
			"/test",
			errors.New("This is a test error"),
		},
		{
			ErrorUnmarshalForm,
			"Unable to unmarshal form data.",
			"/test",
			errors.New("This is a test error"),
		},
		{
			"BadErrorType",
			"",
			"",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(string(test.ErrorType), func(t *testing.T) {
			err := CreateServerError(test.ErrorType, test.URL, test.Err)
			assert.Equal(t, test.Prompt, err.Prompt)
			assert.Equal(t, test.URL, err.URL)
			assert.Equal(t, test.Err, err.Err)
		})
	}

}
