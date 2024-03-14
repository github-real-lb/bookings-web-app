package render

import (
	"bytes"
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
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// TODO: add default templates data
	td.CSRFToken = nosurf.Token(r)
	return td
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
	td = AddDefaultData(td, r)

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

// GetTemplatesCache returns a map of all *.gohtml templates from the directory ./templates
func GetTemplatesCache() (map[string]*template.Template, error) {
	tc := map[string]*template.Template{}

	// get the names of all the files matching *.page.gohtml from ./templates
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return tc, err
	}

	// // get the names of all the files matching *.layout.gohtml from ./templates
	// layouts, err := filepath.Glob("./templates/*.layout.gohtml")
	// if err != nil {
	// 	return tc, err
	// }

	// layoutsExist := len(layouts) > 0

	// range thruogh all the *.page.html files
	for _, page := range pages {
		// extracting the filename itself from the full path
		name := filepath.Base(page)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tc, err
		}

		ts, err = ts.ParseFiles("./templates/base.layout.gohtml")
		if err != nil {
			return tc, err
		}

		if strings.HasSuffix(name, ".room.page.gohtml") {
			ts, err = ts.ParseFiles("./templates/room.layout.gohtml")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
