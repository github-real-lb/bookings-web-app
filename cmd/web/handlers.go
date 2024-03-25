package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/internal/forms"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// HandlerStore provides all handler functions
type HandlerStore struct{}

// NewHandler returns a new http.Handler with all routings
func NewHandler(store HandlerStore) http.Handler {
	mux := chi.NewRouter()

	// add middleware that recover from panics
	mux.Use(middleware.Recoverer)

	// add middleware that loads and saves and session on every request
	mux.Use(app.Session.LoadAndSave)

	// add middleware that provides CSRF protection to all POST requests
	if !app.InTestingMode() {
		mux.Use(NoSurf)
	}

	// setting routes
	mux.Get("/", store.Home)
	mux.Get("/about", store.About)

	mux.Get("/generals-quarters", store.Generals)
	mux.Post("/generals-quarters", store.PostAvailability)

	mux.Get("/majors-suite", store.Majors)
	mux.Post("/majors-suite", store.PostAvailability)

	mux.Get("/contact", store.Contact)

	mux.Get("/search-availability", store.Availability)
	mux.Post("/search-availability", store.PostAvailability)
	mux.Post("/search-availability-json", store.PostAvailabilityJSON)

	mux.Get("/make-reservation", store.Reservation)
	mux.Post("/make-reservation", store.PostReservation)
	mux.Get("/reservation-summary", store.ReservationSummary)

	// setting file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// Home is the home page handler
func (s *HandlerStore) Home(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "home.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// About is the about page handler
func (s *HandlerStore) About(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "about.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Reservation is the generals-quarters room page handler
func (s *HandlerStore) Generals(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "generals.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Majors is the majors-suite room page handler
func (s *HandlerStore) Majors(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "majors.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Contact is the contact page handler
func (s *HandlerStore) Contact(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "contact.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Availability is the search-availability page handler
func (s *HandlerStore) Availability(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostAvailability is the search-availability page post handler
func (s *HandlerStore) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	startDate := r.Form.Get("start_date")
	endDate := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", startDate, endDate)))
}

type jsonResponse struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// AvailabilityJSON handles requests for availability and sends JSON response
func (s *HandlerStore) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", r.Form.Get("start_date"))
	if err != nil {
		app.LogServerError(w, err)
	}
	endDate, err := time.Parse("2006-01-02", r.Form.Get("end_date"))
	if err != nil {
		app.LogServerError(w, err)
	}

	resp := jsonResponse{
		StartDate: startDate,
		EndDate:   endDate,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		app.LogServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation is the make-reservation page handler
func (s *HandlerStore) Reservation(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
		Form: forms.New(nil),
		Data: map[string]any{
			"reservation": Reservation{},
		},
	})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostReservation is the make-reservation post page handler
func (s *HandlerStore) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	reservation := Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLenght("first_name", 3)
	form.MinLenght("last_name", 3)
	form.IsEmailValid("email")

	if !form.Valid() {
		err := RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
			Form: form,
			Data: map[string]any{
				"reservation": reservation,
			},
		})

		if err != nil {
			app.LogServerError(w, err)
		}

		return
	}

	app.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (s *HandlerStore) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		app.LogError(errors.New("cannot get reservation from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.Session.Remove(r.Context(), "reservation")

	err := RenderTemplate(w, r, "reservation-summary.page.gohtml", &TemplateData{
		Data: map[string]any{
			"reservation": reservation,
		},
	})

	if err != nil {
		app.LogServerError(w, err)
	}
}