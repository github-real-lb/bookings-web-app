package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/github-real-lb/bookings-web-app/internal/config"
	"github.com/github-real-lb/bookings-web-app/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplatesCache initiates the gohtml templates cache for the render package
func NewTemplatesCache(ac *config.AppConfig) {
	app = ac
}

// AddDefaultData is used to add default data relevant to all gohtml templates
func AddDefaultData(td *models.TemplateData, r *http.Request) {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
}

// RenderTemplate execute a gohtml template from the template cache.
// It requires to initally assign a template cache using the NewTemplatesCache function.
func RenderTemplate(w http.ResponseWriter, r *http.Request, gohtml string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error

	// UseCache is false in developement mode in order to allow changes of gohtml templates on runtime.
	if app.UseTemplateCache {
		tc = app.TemplateCache
	} else {
		tc, err = GetTemplatesCache()
		if err != nil {
			log.Println(err)
		}
	}

	// checks if gohtml template exist in cache
	t, ok := tc[gohtml]
	if !ok {
		log.Printf("couldn't find %s in template cache.\n", gohtml)
	}

	// adds default templates data relevant to all templates
	AddDefaultData(td, r)

	// check for error in template execution before passing it to w (http.ResponseWriter)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template to w (http.ResponseWriter)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// GetTemplatesCache returns a map of all *.gohtml templates from the directory
// set in AppConfig.TemplatePath
func GetTemplatesCache() (map[string]*template.Template, error) {
	tc := map[string]*template.Template{}

	app.TemplatePath = strings.TrimSuffix(app.TemplatePath, "/")
	pattern := fmt.Sprintf("%s/*.page.gohtml", app.TemplatePath)
	baseFilename := fmt.Sprintf("%s/base.layout.gohtml", app.TemplatePath)
	roomFilename := fmt.Sprintf("%s/room.layout.gohtml", app.TemplatePath)

	// get the names of all the files matching *.page.gohtml from ./templates
	pages, err := filepath.Glob(pattern)
	if err != nil {
		return tc, err
	}

	// range thruogh all the *.page.html files
	for _, page := range pages {
		// extracting the filename itself from the full path
		name := filepath.Base(page)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tc, err
		}

		ts, err = ts.ParseFiles(baseFilename)
		if err != nil {
			return tc, err
		}

		if strings.HasSuffix(name, ".room.page.gohtml") {
			ts, err = ts.ParseFiles(roomFilename)
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
