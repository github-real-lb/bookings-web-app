package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRenderer(t *testing.T) {
	hr := NewRenderer()
	assert.IsType(t, &GoHtmlRenderer{}, hr)
}

func TestGoHtmlRenderer_LoadGoHtmlPageTemplates(t *testing.T) {
	hr := NewRenderer()
	err := hr.LoadGoHtmlPageTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, hr.Templates)
}

func TestGoHtmlRenderer_RenderGoHtmlPageTemplate(t *testing.T) {
	// create new renderer and load templates
	hr := NewRenderer()
	err := hr.LoadGoHtmlPageTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, hr.Templates)

	// create a new test server, and a new request and recorder
	ts := NewTestServer(t)
	request := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	// test ok on reloading templates cache (developement mode)
	app.SetDevelopementMode()
	err = hr.RenderGoHtmlPageTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test ok on using template cache (production modes)
	app.SetTestingMode()
	err = hr.RenderGoHtmlPageTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test not ok on missing template
	err = hr.RenderGoHtmlPageTemplate(recorder, request, "non-existing.page.gohtml", &TemplateData{})
	assert.Error(t, err)
}

func TestAddDefaultData(t *testing.T) {
	// create a new test server and a new request
	ts := NewTestServer(t)
	request := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)

	// add data into session
	app.Session.Put(request.Context(), "flash", "flash")
	app.Session.Put(request.Context(), "warning", "warning")
	app.Session.Put(request.Context(), "error", "error")

	td := TemplateData{}
	addDefaultData(&td, request)
	require.NotEmpty(t, td)
	assert.Equal(t, "flash", td.Flash)
	assert.Equal(t, "warning", td.Warning)
	assert.Equal(t, "error", td.Error)
}

func TestGoHtmlRenderer_LoadGoHtmlMailTemplates(t *testing.T) {
	hr := NewRenderer()
	err := hr.LoadGoHtmlMailTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, hr.Templates)
}

func TestGoHtmlRenderer_RenderGoHtmlMailTemplate(t *testing.T) {
	// create new renderer and load templates
	hr := NewRenderer()
	err := hr.LoadGoHtmlMailTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, hr.Templates)

	// test ok on reloading templates cache (developement and testing modes)
	app.SetDevelopementMode()
	s, err := hr.RenderGoHtmlMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{})
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	// test ok on using template cache (production modes)
	app.SetTestingMode()
	s, err = hr.RenderGoHtmlMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{})
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	// test not ok on missing template
	s, err = hr.RenderGoHtmlMailTemplate("non-existing.page.gohtml", &TemplateData{})
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestGoHtmlRenderer_CreateReservationConfirmationMail(t *testing.T) {
	// create new renderer and load templates
	hr := NewRenderer()
	err := hr.LoadGoHtmlMailTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, hr.Templates)

	// create random reservation
	r := randomReservation()

	mailData, err := hr.CreateReservationConfirmationMail(r)
	require.NoError(t, err)
	assert.Equal(t, r.Email, mailData.To)
	assert.Equal(t, app.Listing.Email, mailData.From)
	assert.Equal(t, fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code), mailData.Subject)
	assert.NotEmpty(t, mailData.Content)
}
