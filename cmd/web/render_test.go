package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddDefaultData(t *testing.T) {
	// create a new test server and a new request
	ts, _ := NewTestServer(t)
	request := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)

	// add data into session
	app.Session.Put(request.Context(), "flash", "flash")
	app.Session.Put(request.Context(), "warning", "warning")
	app.Session.Put(request.Context(), "error", "error")

	td := TemplateData{}
	AddDefaultData(&td, request)
	require.NotEmpty(t, td)
	assert.Equal(t, "flash", td.Flash)
	assert.Equal(t, "warning", td.Warning)
	assert.Equal(t, "error", td.Error)
}

func TestRenderTemplate(t *testing.T) {
	// create a new test server, and a new request and recorder
	ts, _ := NewTestServer(t)
	request := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	// test ok on reloading templates cache (developement and testing modes)
	err := RenderTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test ok on using template cache (production modes)
	app.UseTemplateCache = true
	err = RenderTemplate(recorder, request, "home.page.gohtml", &TemplateData{})
	assert.NoError(t, err)

	// test not ok on missing template
	app.UseTemplateCache = true
	err = RenderTemplate(recorder, request, "non-existing.page.gohtml", &TemplateData{})
	assert.Error(t, err)
}

func TestGetTemplatesCache(t *testing.T) {
	app.SetTestingMode()
	tc, err := GetTemplatesCache()
	require.NoError(t, err)
	assert.NotEmpty(t, tc)
}
