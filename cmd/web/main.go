package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/github-real-lb/go-web-app/pkg/config"
	"github.com/github-real-lb/go-web-app/pkg/handlers"
	"github.com/github-real-lb/go-web-app/pkg/render"
)

const addr = "localhost:8080"

func main() {
	var app config.AppConfig
	var err error

	// set AppConfig to ignore templates cache and loads templates from disk - developement mode
	app.UseCache = false

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
