package main

import (
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util/config"
)

func TestMain(m *testing.M) {
	InitializeApp(config.TestingMode)

	// run tests
	code := m.Run()

	os.Exit(code)
}
