package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testHandler is a mock handler
type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestNoSurf(t *testing.T) {
	h := NoSurf(&testHandler{})
	assert.Implements(t, (*http.Handler)(nil), h)
}

func TestServer_LogRequestsAndResponse(t *testing.T) {
	ts := NewTestServer(t)
	h := ts.LogRequestsAndResponse(&testHandler{})
	assert.Implements(t, (*http.Handler)(nil), h)

	req := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	ts.BuildLogAnyInfoStub()

	h.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuth(t *testing.T) {
	h := Auth(&testHandler{})
	assert.Implements(t, (*http.Handler)(nil), h)

	ts := NewTestServer(t)

	t.Run("Not Authenticated", func(t *testing.T) {
		req := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		h.ServeHTTP(recorder, req)
		v, ok := app.Session.Pop(req.Context(), "error").(string)
		assert.True(t, ok)
		assert.Equal(t, "Access denied. Pleasae log in first!", v)
		assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	})

	t.Run("Authenticated", func(t *testing.T) {
		req := ts.NewRequestWithSession(t, http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		app.Session.Put(req.Context(), "user_id", 1)

		h.ServeHTTP(recorder, req)
		ok := app.Session.Exists(req.Context(), "error")
		assert.False(t, ok)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

}
