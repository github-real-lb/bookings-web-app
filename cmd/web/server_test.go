package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	server := NewServer(ADDRESS)
	assert.IsType(t, (*Server)(nil), server)
}
