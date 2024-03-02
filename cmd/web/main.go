package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/pkg/config"
	"github.com/github-real-lb/bookings-web-app/pkg/handlers"
	"github.com/github-real-lb/bookings-web-app/pkg/render"
)

const addr = "localhost:8080"

var app config.AppConfig

func main() {
	var err error

	// set app to developement mode
	app.InProduction = false

	// setting up session manager
	session := scs.New()
	session.Lifetime = 24 * time.Hour // keeps session data for 24 hours
	session.Cookie.Persist = true     // keeps session data after browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // determines use of SSL encryption
	app.Session = session

	// set app to ignore templates cache and loads templates from disk - developement mode
	app.UseTemplateCache = false

	// load templates cache to AppConfig
	app.TemplateCache, err = render.GetTemplatesCache()
	if err != nil {
		log.Fatal("error creating gohtml templates cache:", err)
	}

	// Initiating the handlers package repo
	handlers.NewHandlersRepository(&app)

	// Initiating the render package templates cahce
	render.NewTemplatesCache(&app)

	// create new http server with routes
	srv := http.Server{
		Addr:    addr,
		Handler: routes(&app),
	}

	// start the server
	fmt.Println("Starting web server on,", addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
