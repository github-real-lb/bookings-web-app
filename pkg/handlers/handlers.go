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

// Reservation is the make-reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Reservation is the generals-quarters room page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.gohtml", &models.TemplateData{})
}

// Majors is the majors-suite room page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.gohtml", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.gohtml", &models.TemplateData{})
}

// Availability is the search-availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.gohtml", &models.TemplateData{})
}
