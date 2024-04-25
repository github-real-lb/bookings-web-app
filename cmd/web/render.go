package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
	"github.com/justinas/nosurf"
)

type TemplateRenderer struct {
	Templates        map[string]*template.Template
	TemplatesPath    string // Path for templates main folder
	UseTemplateCache bool   // Use the loaded templates if true or reload from disk if false
}

// NewTemplateRenderer creates a new empty instance of TemplateRenderer.
// path is the templates directory path
func NewTemplateRenderer(path string) *TemplateRenderer {
	path = strings.TrimSuffix(path, "/")
	return &TemplateRenderer{
		Templates:        make(map[string]*template.Template),
		TemplatesPath:    path,
		UseTemplateCache: false,
	}
}

// LoadGoTemplates loads all templates from TemplateCache.Path
func (tr *TemplateRenderer) LoadGoTemplates() error {
	if tr.TemplatesPath == "" || tr.Templates == nil {
		return fmt.Errorf("templateCache needs to be initiated with NewTemplateCache function")
	}

	// Load site pages gohtml templates
	dir := "pages"
	pattern := "*.page.gohtml"

	err := tr.loadGoTemplatesFromDirectory(dir, pattern)
	if err != nil {
		return err
	}

	// Load admin pages gohtml templates
	dir = "admin"
	pattern = "*.panel.gohtml"

	err = tr.loadGoTemplatesFromDirectory(dir, pattern)
	if err != nil {
		return err
	}

	// Load mail gohtml templates
	dir = "mails"
	pattern = "*.mail.gohtml"

	err = tr.loadGoTemplatesFromDirectory(dir, pattern)
	if err != nil {
		return err
	}

	return nil
}

// loadGoTemplatesFromDirectory loads specific templates from specific sub-directory in TemplateCache.Path
func (tr *TemplateRenderer) loadGoTemplatesFromDirectory(dir string, pattern string) error {
	pattern = fmt.Sprintf("%s/%s/%s", tr.TemplatesPath, dir, pattern)
	baseFilename := fmt.Sprintf("%s/%s/base.layout.gohtml", tr.TemplatesPath, dir)

	// get the names of all the files matching pagePattern
	pages, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	// range thruogh all the pagePattern files
	for _, page := range pages {
		// extracting the filename itself from the full path
		name := filepath.Base(page)

		// creating a new template set with the page name, and parsing the gohtml page.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return err
		}

		ts, err = ts.ParseFiles(baseFilename)
		if err != nil {
			return err
		}

		tr.Templates[name] = ts
	}

	return nil
}

// RenderGoTemplate execute a gohtml template
func (tr *TemplateRenderer) RenderGoTemplate(w http.ResponseWriter, r *http.Request, gohtml string, td *TemplateData) error {
	if tr.TemplatesPath == "" || tr.Templates == nil {
		return fmt.Errorf("templateCache needs to be initiated with NewTemplateCache function")
	}

	var err error

	// Load Templates from disk in developement mode in order to allow template updates on runtime.
	if !tr.UseTemplateCache {
		err = tr.LoadGoTemplates()
		if err != nil {
			return err
		}
	}

	// checks if gohtml template exist in cache
	t, ok := tr.Templates[gohtml]
	if !ok {
		return fmt.Errorf("couldn't find %s in template cache", gohtml)
	}

	// adds default templates data relevant to all templates
	addDefaultData(td, r)

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

// AddDefaultData is used to add default data relevant to all gohtml templates
func addDefaultData(td *TemplateData, r *http.Request) {
	// add green success messages
	if app.Session.Exists(r.Context(), "flash") {
		td.Flash = app.Session.PopString(r.Context(), "flash")
	}

	// add yellow warning messages
	if app.Session.Exists(r.Context(), "warning") {
		td.Warning = app.Session.PopString(r.Context(), "warning")
	}

	// add red error messages
	if app.Session.Exists(r.Context(), "error") {
		td.Error = app.Session.PopString(r.Context(), "error")
	}

	// set login status
	td.IsAuthenticated = IsAuthenticated(r)

	// add listing information
	td.Listing = app.Listing

	// add CSRF Token
	td.CSRFToken = nosurf.Token(r)
}

// RenderMailTemplate execute an email gohtml template
func (tr *TemplateRenderer) RenderMailTemplate(gohtml string, td *TemplateData) (string, error) {
	if tr.TemplatesPath == "" || tr.Templates == nil {
		return "", fmt.Errorf("templateCache needs to be initiated with NewTemplateCache function")
	}

	var err error

	// Load Templates from disk in developement mode in order to allow template updates on runtime.
	if !tr.UseTemplateCache {
		err = tr.LoadGoTemplates()
		if err != nil {
			return "", err
		}
	}

	// checks if gohtml template exist in cache
	t, ok := tr.Templates[gohtml]
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

// CreateReservationNotificationMail creates reservation confirmation mail
func (tr *TemplateRenderer) CreateReservationConfirmationMail(r Reservation) (mailers.MailData, error) {
	var err error

	// create reservation confirmation email
	data := mailers.MailData{
		To:      r.Email,
		From:    app.Listing.Email,
		Subject: fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code),
	}

	data.Content, err = tr.RenderMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{
		Data: map[string]any{
			"start_date":  r.StartDate.Format(config.DateLayout),
			"end_date":    r.EndDate.Format(config.DateLayout),
			"reservation": r,
		},
	})

	return data, err
}
