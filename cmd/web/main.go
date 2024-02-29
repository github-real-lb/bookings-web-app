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

	// load
	app.TemplateCache, err = render.GetTemplatesCache()
	if err != nil {
		log.Fatal("error creating gohtml templates cache:", err)
	}

	render.LoadTemplatesCache(&app)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println("Starting web server on,", addr)
	http.ListenAndServe(addr, nil)
}
