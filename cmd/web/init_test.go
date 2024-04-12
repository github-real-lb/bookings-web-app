package main

import (
	"testing"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeApp(t *testing.T) {
	t.Run("TestingMode", func(t *testing.T) {
		err := InitializeApp(config.TestingMode)
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.AppConfig)
		require.NotNil(t, app.DBConfig)
		assert.Equal(t, config.TestingMode, app.AppConfig.Mode)

	})
}
