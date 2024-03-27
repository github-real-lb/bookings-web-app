package db

import (
	"os"
	"testing"

	"github.com/github-real-lb/bookings-web-app/util/config"
)

const DBConfigFilename = "./../db.config.json"

var dbConfig config.DBConfig

func TestMain(m *testing.M) {
	var err error
	dbConfig, err = config.LoadDBConfig(DBConfigFilename)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
