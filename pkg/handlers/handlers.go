package handlers

import (
	"net/http"

	"github.com/github-real-lb/bookings-web-app/models"
	"github.com/github-real-lb/bookings-web-app/pkg/config"
	"github.com/github-real-lb/bookings-web-app/pkg/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by the handlers
var Repo *Repository

// // NewRepo creates a new Repository of Handlers
// func NewRepo(ac *config.AppConfig) *Repository {
// 	return &Repository{
// 		App: ac,
// 	}
// }

// NewHandlers initiates the repository for the handlers package
func NewHandlersRepository(ac *config.AppConfig) {
	Repo = &Repository{
		App: ac,
	}
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
