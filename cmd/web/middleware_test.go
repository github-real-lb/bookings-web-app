package main

import (
	"net/http"
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
}
