package web

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// AppMode defines the application modes: Production = 0, Development = 1, Testing = 2
type AppMode int

const (
	ProductionMode  AppMode = 0
	DevelopmentMode AppMode = 1
	TestingMode     AppMode = 2
)

// AppConfig holds the application config
type AppConfig struct {
	AppMode

	AppLogger

	// ServerAddress is http.Server listening address
	ServerAddress string

	// Session is the session manager
	Session *scs.SessionManager

	// TemplatePath is the path to the templates folder.
	// Default: "./templates"
	TemplatePath string

	// TemplateCache is a memory cache for all gohtml pages
	TemplateCache map[string]*template.Template

	// UseTemplateCache determines if gohtml pages load from cache (true) or disk (false).
	// use false in developement mode.
	UseTemplateCache bool
}

// LoadConfig returns the Application Configuration set for InProductionMode.
func LoadConfig() AppConfig {
	var app AppConfig

	// setting app configuration to production mode
	app.AppMode = ProductionMode
	app.UseTemplateCache = true

	// initializing loggers
	app.AppLogger = NewAppLogger()

	// TODO: these should be loaded from app.env file.
	app.ServerAddress = "localhost:8080"
	app.TemplatePath = "./templates"

	return app
}

// SetDevelopementMode sets the Application Configuration for development.
func (app *AppConfig) SetDevelopementMode() {
	app.AppMode = DevelopmentMode
	app.UseTemplateCache = false
}

// SetTestingMode sets the Application Configuration for testing.
func (app *AppConfig) SetTestingMode() {
	app.AppMode = TestingMode
	app.UseTemplateCache = false

	// TODO: this should be loaded from app.env file.
	app.TemplatePath = "./../../templates"
}

// InProductionMode returns true if the Application Configuration is set for production.
func (app *AppConfig) InProductionMode() bool {
	return app.AppMode == ProductionMode
}

// InDevelopmentMode returns true if the Application Configuration is set for development.
func (app *AppConfig) InDevelopmentMode() bool {
	return app.AppMode == DevelopmentMode
}

// InTestingMode returns true if the Application Configuration is set for testing.
func (app *AppConfig) InTestingMode() bool {
	return app.AppMode == TestingMode
}
