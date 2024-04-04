package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/github-real-lb/bookings-web-app/util/loggers"
)

// AppMode defines the application modes: Production = 0, Development = 1, Testing = 2
type AppMode int

const (
	ProductionMode  AppMode = 0
	DevelopmentMode AppMode = 1
	TestingMode     AppMode = 2
	DebuggingMode   AppMode = 3
)

const DateTimeLayout = "2006-01-02 15:04:05.999999999Z07:00"
const DateLayout = "2006-01-02"

// AppConfig holds the application config
type AppConfig struct {
	AppMode

	loggers.AppLogger

	// Session is the session manager
	Session *scs.SessionManager

	// StartingPathProduction is the production starting path of the app.
	StartingPathProduction string `json:"starting_path_production"`

	// StartingPathTesting is the testing starting path of the app.
	StartingPathTesting string `json:"starting_path_testing"`

	// StaticDirectoryName is the directory of the static files of the app.
	StaticDirectoryName string `json:"static_directory_name"`

	// StaticPath is the full path of the static folder.
	StaticPath string

	// TemplateCache is a memory cache for all gohtml pages.
	TemplateCache map[string]*template.Template

	// TemplateDirectoryName is the name of the templates folder.
	TemplateDirectoryName string `json:"template_directory_name"`

	// TemplatePath is the full path of the templates folder.
	TemplatePath string

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

	app.TemplateDirectoryName = strings.TrimSuffix(app.TemplateDirectoryName, "/")
	app.StaticDirectoryName = strings.TrimSuffix(app.StaticDirectoryName, "/")

	// initializing loggers
	app.AppLogger = loggers.NewAppLogger()

	switch appMode {
	case ProductionMode:
		app.SetProductionMode()
	case DevelopmentMode:
		app.SetDevelopementMode()
	case TestingMode:
		app.SetTestingMode()
	case DebuggingMode:
		app.SetDebuggingMode()
	default:
		return &app, errors.New("invalid application mode setting in config file")
	}

	return &app, nil
}

// SetProductionMode sets the Application Configuration for production.
func (app *AppConfig) SetProductionMode() {
	app.AppMode = ProductionMode
	app.TemplatePath = fmt.Sprint(app.StartingPathProduction, app.TemplateDirectoryName)
	app.StaticPath = fmt.Sprint(app.StartingPathProduction, app.StaticDirectoryName)
	app.UseTemplateCache = true
}

// SetDevelopementMode sets the Application Configuration for development.
func (app *AppConfig) SetDevelopementMode() {
	app.AppMode = DevelopmentMode
	app.TemplatePath = fmt.Sprint(app.StartingPathProduction, app.TemplateDirectoryName)
	app.StaticPath = fmt.Sprint(app.StartingPathProduction, app.StaticDirectoryName)
	app.UseTemplateCache = false
}

// SetTestingMode sets the Application Configuration for testing.
func (app *AppConfig) SetTestingMode() {
	app.AppMode = TestingMode
	app.TemplatePath = fmt.Sprint(app.StartingPathTesting, app.TemplateDirectoryName)
	app.StaticPath = fmt.Sprint(app.StartingPathTesting, app.StaticDirectoryName)
	app.UseTemplateCache = false
}

// SetDebuggingMode sets the Application Configuration for debugging with the IDE debugger.
func (app *AppConfig) SetDebuggingMode() {
	app.AppMode = DebuggingMode
	app.TemplatePath = fmt.Sprint(app.StartingPathTesting, app.TemplateDirectoryName)
	app.StaticPath = fmt.Sprint(app.StartingPathTesting, app.StaticDirectoryName)
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

// InTestingMode returns true if the Application Configuration is set for testing.
func (app *AppConfig) InDebuggingMode() bool {
	return app.AppMode == DebuggingMode
}
