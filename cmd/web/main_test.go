package main

import (
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestInitializeApp(t *testing.T) {
	t.Run("ProductionMode", func(t *testing.T) {
		err := InitializeApp(config.ProductionMode)
		assert.NoError(t, err)
	})

	t.Run("DevelopmentMode", func(t *testing.T) {
		err := InitializeApp(config.DevelopmentMode)
		assert.NoError(t, err)
	})

	t.Run("TestingMode", func(t *testing.T) {
		err := InitializeApp(config.TestingMode)
		assert.NoError(t, err)
	})

}
