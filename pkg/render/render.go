package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/github-real-lb/go-web-app/pkg/config"
)

var app *config.AppConfig

func LoadTemplatesCache(ac *config.AppConfig) {
	app = ac
}

// RenderTemplate execute a gohtml template from the template cache.
// It requires that the templates cache will be loaded initally using the LoadTemplatesCache function.
func RenderTemplate(w http.ResponseWriter, gohtml string) {
	t, ok := app.TemplateCache[gohtml]
	if !ok {
		log.Printf("couldn't find %s in template cache\n", gohtml)
	}

	// check for error in template execution before passing it to w (http.ResponseWriter)
	buf := new(bytes.Buffer)
	err := t.Execute(buf, nil)
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

	// get the names of all the files matching *.layout.gohtml from ./templates
	layouts, err := filepath.Glob("./templates/*.layout.gohtml")
	if err != nil {
		return tc, err
	}

	layoutsExist := len(layouts) > 0

	// range thruogh all the *.page.html files
	for _, page := range pages {
		// extracting the filename itself from the full path
		name := filepath.Base(page)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tc, err
		}

		if layoutsExist {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
