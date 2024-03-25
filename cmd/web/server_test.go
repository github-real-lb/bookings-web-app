package main

import (
	"testing"

	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	store := mocks.NewMockStore(t)

	server := NewServer(store)
	assert.IsType(t, (*Server)(nil), server)
}
