package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/util/config"
)

const (
	AppConfigFilename = "./../../app.config.json"
	DBConfigFilename  = "./../../db.config.json"
)

type AppConfig struct {
	config.AppConfig
	config.DBConfig
}

// app holds the configurations and templates of the app
var app AppConfig

// InitializeApp loads the app configurations and setup based on the application mode
func InitializeApp(appMode config.AppMode) error {
	var err error

	// load application default configurations
	app.AppConfig, err = config.LoadAppConfig(AppConfigFilename, appMode)
	if err != nil {
		return err
	}

	app.DBConfig, err = config.LoadDBConfig(DBConfigFilename)
	if err != nil {
		return err
	}

	// load templates cache to AppConfig
	app.TemplateCache, err = GetTemplatesCache()
	if err != nil {
		return errors.New(fmt.Sprint("error creating gohtml templates cache: ", err.Error()))
	}

	// setting up session manager
	session := scs.New()
	session.Lifetime = 24 * time.Hour // keeps session data for 24 hours
	session.Cookie.Persist = true     // keeps session data after browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProductionMode() // determines use of SSL encryption
	app.Session = session

	// defining session stored types
	gob.Register(User{})
	gob.Register(Room{})
	gob.Register(Restriction{})
	gob.Register(Reservation{})
	gob.Register(RoomRestriction{})

	return nil
}
