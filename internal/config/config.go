package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	// InProduction determines if the webapp is running in production mode (true) or developement mode (false)
	InProduction bool

	// UseTemplateCache determines if gohtml pages load from cache (true) or disk (false).
	// use false in developement mode.
	UseTemplateCache bool

	// TemplatePath is the path to the templates folder.
	// Default: "./templates"
	TemplatePath string

	// TemplateCache is a memory cache for all gohtml pages
	TemplateCache map[string]*template.Template

	// Session is the session manager
	Session *scs.SessionManager

	InfoLog *log.Logger
}

func (app *AppConfig) InDevelopementMode() {
	app.InProduction = false
	app.UseTemplateCache = false
}
