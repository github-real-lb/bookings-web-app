package main

import (
	"testing"

	"github.com/github-real-lb/bookings-web-app/db/mocks"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	store := mocks.NewMockStore(t)

	server := NewServer(store)
	require.IsType(t, (*Server)(nil), server)
}
