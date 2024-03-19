package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/internal/config"
	"github.com/github-real-lb/bookings-web-app/internal/models"
	"github.com/github-real-lb/bookings-web-app/internal/render"
)

const ADDRESS = "localhost:8080"

var templatePath = "./templates"
var app config.AppConfig

func main() {
	err := InitApp()
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(ADDRESS)

	// start the server
	fmt.Println("Starting web server on,", ADDRESS)
	err = server.ListenAndServe()
	log.Fatal(err)
}

func InitApp() error {
	var err error

	// set app to developement mode
	app.InDevelopementMode()

	// Initiating the render package templates cahce
	render.NewTemplatesCache(&app)

	// load templates cache to AppConfig
	app.TemplatePath = templatePath
	app.TemplateCache, err = render.GetTemplatesCache()
	if err != nil {
		return errors.New(fmt.Sprint("error creating gohtml templates cache: ", err.Error()))
	}

	// setting up session manager
	session := scs.New()
	session.Lifetime = 24 * time.Hour // keeps session data for 24 hours
	session.Cookie.Persist = true     // keeps session data after browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // determines use of SSL encryption
	app.Session = session

	// defining session stored types
	gob.Register(models.Reservation{})

	return nil
}
