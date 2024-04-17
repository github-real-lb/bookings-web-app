package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

// AddDefaultData is used to add default data relevant to all gohtml templates
func AddDefaultData(td *TemplateData, r *http.Request) {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
}

// GetTemplatesCache returns a map of all *.gohtml templates from the directory
// set in AppConfig.TemplatePath
func GetTemplatesCache() (map[string]*template.Template, error) {
	tc := map[string]*template.Template{}

	pagePattern := fmt.Sprintf("%s/pages/*.page.gohtml", app.TemplatePath)
	mailPattern := fmt.Sprintf("%s/mails/*.mail.gohtml", app.TemplatePath)
	basePageFilename := fmt.Sprintf("%s/pages/base.layout.gohtml", app.TemplatePath)
	baseMailFilename := fmt.Sprintf("%s/mails/base.layout.gohtml", app.TemplatePath)

	// get the names of all the files matching *.page.gohtml from ./templates
	pages, err := filepath.Glob(pagePattern)
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

		ts, err = ts.ParseFiles(basePageFilename)
		if err != nil {
			return tc, err
		}

		tc[name] = ts
	}

	// get the names of all the files matching *.page.gohtml from ./templates
	mailPages, err := filepath.Glob(mailPattern)
	if err != nil {
		return tc, err
	}

	// range thruogh all the *.mail.html files
	for _, page := range mailPages {
		// extracting the filename itself from the full path
		name := filepath.Base(page)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tc, err
		}

		ts, err = ts.ParseFiles(baseMailFilename)
		if err != nil {
			return tc, err
		}

		tc[name] = ts
	}

	return tc, nil
}

// RenderTemplate execute a gohtml template
func RenderTemplate(w http.ResponseWriter, r *http.Request, gohtml string, td *TemplateData) error {
	var tc map[string]*template.Template
	var err error

	// UseCache is false in developement mode in order to allow changes of gohtml templates on runtime.
	if app.UseTemplateCache {
		tc = app.TemplateCache
	} else {
		tc, err = GetTemplatesCache()
		if err != nil {
			return err
		}
	}

	// checks if gohtml template exist in cache
	t, ok := tc[gohtml]
	if !ok {
		return fmt.Errorf("couldn't find %s in template cache", gohtml)
	}

	// adds default templates data relevant to all templates
	AddDefaultData(td, r)

	// check for error in template execution before passing it to w (http.ResponseWriter)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, td)
	if err != nil {
		return err
	}

	// render the template to w (http.ResponseWriter)
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

// RenderMailTemplate execute an email gohtml template
func RenderMailTemplate(gohtml string, td *TemplateData) (string, error) {
	var tc map[string]*template.Template
	var err error

	// UseCache is false in developement mode in order to allow changes of gohtml templates on runtime.
	if app.UseTemplateCache {
		tc = app.TemplateCache
	} else {
		tc, err = GetTemplatesCache()
		if err != nil {
			return "", err
		}
	}

	// checks if gohtml template exist in cache
	t, ok := tc[gohtml]
	if !ok {
		return "", fmt.Errorf("couldn't find %s in template cache", gohtml)
	}

	// check for error in template execution before passing it to w (http.ResponseWriter)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, td)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
