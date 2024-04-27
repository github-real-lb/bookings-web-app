package main

import (
	"fmt"
	"net/http"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
	"github.com/github-real-lb/bookings-web-app/util/render"
	"github.com/justinas/nosurf"
)

type GoHtmlRenderer struct {
	*render.SmartRenderer
}

func NewRenderer() *GoHtmlRenderer {
	return &GoHtmlRenderer{
		SmartRenderer: render.NewSmartRenderer(),
	}
}

// LoadGoHtmlPageTemplates loads all web page templates
func (hr *GoHtmlRenderer) LoadGoHtmlPageTemplates() error {
	// Load site pages gohtml templates
	path := fmt.Sprintf("%s/%s", app.TemplatePath, "pages")
	basePattern := "base.layout.gohtml"
	tmplPattern := "*.page.gohtml"

	err := hr.LoadTemplates(path, basePattern, tmplPattern)
	if err != nil {
		return err
	}

	// Load admin pages gohtml templates
	path = fmt.Sprintf("%s/%s", app.TemplatePath, "admin")
	basePattern = "base.layout.gohtml"
	tmplPattern = "*.panel.gohtml"

	err = hr.LoadTemplates(path, basePattern, tmplPattern)
	if err != nil {
		return err
	}

	return nil
}

// RenderGoHtmlPageTemplate executes a web page gohtml template
func (hr *GoHtmlRenderer) RenderGoHtmlPageTemplate(w http.ResponseWriter, r *http.Request, gohtml string, td *TemplateData) error {
	var err error

	// load Templates from disk in developement mode in order to allow template updates on runtime.
	if app.InDevelopmentMode() {
		err = hr.LoadGoHtmlPageTemplates()
		if err != nil {
			return err
		}
	}

	// add default templates data relevant to all templates
	addDefaultData(td, r)

	// render template
	data, err := hr.RenderTemplate(gohtml, td)
	if err != nil {
		return err
	}

	// render the template to w (http.ResponseWriter)
	_, err = w.Write(data)
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

// LoadGoHtmlMailTemplates loads all email templates
func (hr *GoHtmlRenderer) LoadGoHtmlMailTemplates() error {
	// Load mail gohtml templates
	path := fmt.Sprintf("%s/%s", app.TemplatePath, "mails")
	basePattern := "base.layout.gohtml"
	tmplPattern := "*.mail.gohtml"

	err := hr.LoadTemplates(path, basePattern, tmplPattern)
	if err != nil {
		return err
	}

	return nil
}

// RenderGoHtmlMailTemplate execute an email gohtml template
func (hr *GoHtmlRenderer) RenderGoHtmlMailTemplate(gohtml string, td *TemplateData) (string, error) {
	var err error

	// load Templates from disk in developement mode in order to allow template updates on runtime.
	if app.InDevelopmentMode() {
		err = hr.LoadGoHtmlMailTemplates()
		if err != nil {
			return "", err
		}
	}

	// render template
	data, err := hr.RenderTemplate(gohtml, td)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateReservationNotificationMail creates reservation confirmation mail
func (hr *GoHtmlRenderer) CreateReservationConfirmationMail(r Reservation) (mailers.MailData, error) {
	var err error

	// create reservation confirmation email
	data := mailers.MailData{
		To:      r.Email,
		From:    app.Listing.Email,
		Subject: fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code),
	}

	data.Content, err = hr.RenderGoHtmlMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{
		Data: map[string]any{
			"start_date":  r.StartDate.Format(config.DateLayout),
			"end_date":    r.EndDate.Format(config.DateLayout),
			"reservation": r,
		},
	})

	return data, err
}
