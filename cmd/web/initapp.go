package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/util/webapp"
)

// app holds the configurations and templates of the app
var app webapp.AppConfig

// InitializeApp loads the app configurations and setup based on the application mode
func InitializeApp(appMode webapp.AppMode) error {
	var err error

	// load application default configurations
	app = webapp.LoadConfig()

	// setting up application to required mode
	switch appMode {
	case webapp.DevelopmentMode:
		app.SetDevelopementMode()
	case webapp.TestingMode:
		app.SetTestingMode()
	default:
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
	gob.Register(Reservation{})

	return nil
}
