package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestInitApp(t *testing.T) {
	err := InitApp()
	assert.NoError(t, err)
}
