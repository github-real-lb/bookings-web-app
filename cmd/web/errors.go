package main

import (
	"errors"
	"fmt"

	"github.com/github-real-lb/bookings-web-app/util"
)

type ServerError struct {
	Prompt string
	URL    string
	Err    error
}

func (e ServerError) Error() string {
	text := util.NewText()

	if e.Err != nil {
		text.AddLineIndent(e.Err.Error(), "\t")
	}

	if e.Prompt != "" {
		text.AddLineIndent(fmt.Sprint("PROMPT: ", e.Prompt), "\t")
	}
	if e.URL != "" {
		text.AddLineIndent(fmt.Sprint("URL: ", e.URL), "\t")
	}

	return text.String()
}

type ErrorType string

const (
	ErrorMissingReservation ErrorType = "MissingReservation"
	ErrorParseForm          ErrorType = "ParseForm"
	ErrorRenderTemplate     ErrorType = "RenderTemplate"
	ErrorInvalidParameter   ErrorType = "Invalid Parameter"
)

func CreateServerError(errType ErrorType, url string, err error) ServerError {
	switch errType {
	case ErrorMissingReservation:
		return ServerError{
			Prompt: "No reservation exists. Please make a reservation.",
			URL:    url,
			Err:    errors.New("wrong routing"),
		}
	case ErrorParseForm:
		return ServerError{
			Prompt: "Unable to parse form.",
			URL:    url,
			Err:    err,
		}
	case ErrorRenderTemplate:
		return ServerError{
			Prompt: "Unable to render template.",
			URL:    url,
			Err:    err,
		}
	case ErrorInvalidParameter:
		return ServerError{
			Prompt: "Invalid parameter in URL.",
			URL:    url,
			Err:    errors.New("wrong routing"),
		}
	default:
		return ServerError{}
	}
}
