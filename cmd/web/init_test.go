package main

import (
	"testing"

	"github.com/github-real-lb/bookings-web-app/util/webapp"
	"github.com/stretchr/testify/assert"
)

func TestInitializeApp(t *testing.T) {
	t.Run("ProductionMode", func(t *testing.T) {
		err := InitializeApp(webapp.ProductionMode)
		assert.NoError(t, err)
	})

	t.Run("DevelopmentMode", func(t *testing.T) {
		err := InitializeApp(webapp.DevelopmentMode)
		assert.NoError(t, err)

	})

	t.Run("TestingMode", func(t *testing.T) {
		err := InitializeApp(webapp.TestingMode)
		assert.NoError(t, err)

	})
}
