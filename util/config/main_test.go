package config

import (
	"os"
	"testing"
)

const testAppConfigFilename = "./../../app.config.json"

var app *AppConfig

func TestMain(m *testing.M) {
	var err error
	app, err = LoadAppConfig(testAppConfigFilename, TestingMode)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
