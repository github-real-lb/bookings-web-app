package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateRenderer(t *testing.T) {
	tr := NewTemplateRenderer(app.TemplatePath)
	assert.NotNil(t, tr.Templates)
	assert.Equal(t, strings.TrimSuffix(app.TemplatePath, "/"), tr.TemplatesPath)
	assert.False(t, tr.UseTemplateCache)
}

func TestTemplateRenderer_LoadGoTemplates(t *testing.T) {
	tr := NewTemplateRenderer(app.TemplatePath)
	err := tr.LoadGoTemplates()
	assert.NoError(t, err)
	assert.NotEmpty(t, tr.Templates)
}

func TestTemplateRenderer_loadGoTemplatesFromDirectory(t *testing.T) {
	tests := []struct {
		Name    string
		Dir     string
		Pattern string
		OK      bool
	}{
		{"OK Pages Dir", "pages", "*.page.gohtml", true},
		{"OK Admin Dir", "admin", "*.panel.gohtml", true},
		{"OK Mails Dir", "mails", "*.mail.gohtml", true},
		{"Error", "", "", false},
	}

	for _, test := range tests {
		tr := NewTemplateRenderer(app.TemplatePath)

		err := tr.loadGoTemplatesFromDirectory(test.Dir, test.Pattern)
		if test.OK {
			assert.NoError(t, err)
			assert.NotEmpty(t, tr.Templates)
		} else {
			assert.Error(t, err)
			assert.Empty(t, tr.Templates)
		}
	}
}

func TestTemplateRenderer_RenderGoTemplate(t *testing.T) {
	// create new renderer and load templates
	tr := NewTemplateRenderer(app.TemplatePath)
	err := tr.LoadGoTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, tr.Templates)

	// create a new test server, and a new request and recorder
	ts := NewTestServer(t)
	request := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	// test ok on reloading templates cache (developement and testing modes)
	tr.UseTemplateCache = false
	err = tr.RenderGoTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test ok on using template cache (production modes)
	tr.UseTemplateCache = true
	err = tr.RenderGoTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test not ok on missing template
	err = tr.RenderGoTemplate(recorder, request, "non-existing.page.gohtml", &TemplateData{})
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

func TestTemplateRenderer_RenderMailTemplate(t *testing.T) {
	// create new renderer and load templates
	tr := NewTemplateRenderer(app.TemplatePath)
	err := tr.LoadGoTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, tr.Templates)

	// test ok on reloading templates cache (developement and testing modes)
	tr.UseTemplateCache = false
	s, err := tr.RenderMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{})
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	// test ok on using template cache (production modes)
	tr.UseTemplateCache = true
	s, err = tr.RenderMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{})
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	// test not ok on missing template
	s, err = tr.RenderMailTemplate("non-existing.page.gohtml", &TemplateData{})
	assert.Error(t, err)
	assert.Empty(t, s)
}

func TestTemplateRenderer_CreateReservationConfirmationMail(t *testing.T) {
	// create new renderer and load templates
	tr := NewTemplateRenderer(app.TemplatePath)
	err := tr.LoadGoTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, tr.Templates)

	// create random reservation
	r := randomReservation()

	mailData, err := tr.CreateReservationConfirmationMail(r)
	require.NoError(t, err)
	assert.Equal(t, r.Email, mailData.To)
	assert.Equal(t, app.Listing.Email, mailData.From)
	assert.Equal(t, fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code), mailData.Subject)
	assert.NotEmpty(t, mailData.Content)
}
