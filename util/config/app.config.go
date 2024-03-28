package config

import (
	"encoding/json"
	"errors"
	"html/template"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/util/loggers"
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

	loggers.AppLogger

	// Session is the session manager
	Session *scs.SessionManager

	// TemplatePath is the path to the templates folder.
	TemplatePath string

	// TemplatePath is the path to the templates folder.
	TemplatePathProduction string `json:"template_path_production"`

	// TemplatePath is the path to the templates folder in testing mode.
	TemplatePathTesting string `json:"template_path_testing"`

	// TemplateCache is a memory cache for all gohtml pages
	TemplateCache map[string]*template.Template

	// UseTemplateCache determines if gohtml pages load from cache (true) or disk (false).
	// use false in developement mode.
	UseTemplateCache bool
}

// LoadConfig returns the Application Configuration.
func LoadAppConfig(filename string, appMode AppMode) (*AppConfig, error) {
	app := AppConfig{}

	file, err := os.Open(filename)
	if err != nil {
		return &app, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		return &app, err
	}

	// initializing loggers
	app.AppLogger = loggers.NewAppLogger()

	switch appMode {
	case ProductionMode:
		app.SetProductionMode()
	case DevelopmentMode:
		app.SetDevelopementMode()
	case TestingMode:
		app.SetTestingMode()
	default:
		return &app, errors.New("invalid application mode setting in config file")
	}

	return &app, nil
}

// SetDevelopementMode sets the Application Configuration for development.
func (app *AppConfig) SetProductionMode() {
	app.AppMode = ProductionMode
	app.TemplatePath = app.TemplatePathProduction
	app.UseTemplateCache = true
}

// SetDevelopementMode sets the Application Configuration for development.
func (app *AppConfig) SetDevelopementMode() {
	app.AppMode = DevelopmentMode
	app.TemplatePath = app.TemplatePathProduction
	app.UseTemplateCache = false
}

// SetTestingMode sets the Application Configuration for testing.
func (app *AppConfig) SetTestingMode() {
	app.AppMode = TestingMode
	app.TemplatePath = app.TemplatePathTesting
	app.UseTemplateCache = false
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
