package db

import (
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util/config"
)

const DBConfigFilename = "./../db.config.json"

type StringMap map[string]string

var testStore DatabaseStore

func TestMain(m *testing.M) {
	dbConfig, err := config.LoadDBConfig(DBConfigFilename)
	if err != nil {
		os.Exit(1)
	}

	testStore, err = NewPostgresDBStore(dbConfig.ConnectionString)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
