package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/util/web"
)

// app holds the configurations and templates of the app
var app web.AppConfig

func main() {
	err := InitializeApp(web.DevelopmentMode)
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(app.ServerAddress)

	// start the server
	fmt.Println("Starting web server on,", app.ServerAddress)
	err = server.ListenAndServe()
	log.Fatal(err)
}

// InitializeApp loads the app configurations and setup based on the application mode
func InitializeApp(appMode web.AppMode) error {
	var err error

	// load application default configurations
	app = web.LoadConfig()

	// setting up application to required mode
	switch appMode {
	case web.DevelopmentMode:
		app.SetDevelopementMode()
	case web.TestingMode:
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

func dbtest() {

}
