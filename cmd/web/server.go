package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/internal/forms"
	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server handles all routing and provides all database functions
type Server struct {
	Router        *http.Server
	DatabaseStore db.DatabaseStore
}

// NewServer returns a new Server with Router and Database Store
func NewServer(store db.DatabaseStore) *Server {
	mux := chi.NewRouter()

	server := Server{
		Router: &http.Server{
			Addr:    app.ServerAddress,
			Handler: mux,
		},
		DatabaseStore: store,
	}

	// add middleware that recover from panics
	mux.Use(middleware.Recoverer)

	// add middleware that loads and saves and session on every request
	mux.Use(app.Session.LoadAndSave)

	// add middleware that provides CSRF protection to all POST requests
	if !app.InTestingMode() {
		mux.Use(NoSurf)
	}

	// setting routes
	mux.Get("/", server.Home)
	mux.Get("/about", server.About)

	mux.Get("/generals-quarters", server.Generals)
	mux.Post("/generals-quarters", server.PostAvailability)

	mux.Get("/majors-suite", server.Majors)
	mux.Post("/majors-suite", server.PostAvailability)

	mux.Get("/contact", server.Contact)

	mux.Get("/search-availability", server.Availability)
	mux.Post("/search-availability", server.PostAvailability)
	mux.Post("/search-availability-json", server.PostAvailabilityJSON)

	mux.Get("/make-reservation", server.Reservation)
	mux.Post("/make-reservation", server.PostReservation)
	mux.Get("/reservation-summary", server.ReservationSummary)

	// setting file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return &server
}

// Home is the home page handler
func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "home.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// About is the about page handler
func (s *Server) About(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "about.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Reservation is the generals-quarters room page handler
func (s *Server) Generals(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "generals.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Majors is the majors-suite room page handler
func (s *Server) Majors(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "majors.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Contact is the contact page handler
func (s *Server) Contact(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "contact.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// Availability is the search-availability page handler
func (s *Server) Availability(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostAvailability is the search-availability page post handler
func (s *Server) PostAvailability(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) Reservation(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	var form *forms.Form
	var reservation Reservation

	// create a new form with data and validate the form
	form = forms.New(r.PostForm)

	//TODO replace this code part with code loading the data from the session
	form.Add("start_date", time.Now().Format("2006-01-02"))
	form.Add("end_date", time.Now().Add(time.Hour*24*7).Format("2006-01-02"))
	form.Add("room_id", "1")
	//TODO end

	form.TrimSpaces()
	form.Required("first_name", "last_name", "email")
	form.MinLenght("first_name", 3)
	form.MinLenght("last_name", 3)
	form.IsEmailValid("email")

	// Parse form's data to reservation
	reservationData := form.Marshal()
	err = reservation.Unmarshal(reservationData)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

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

	// define create reservation parameters
	arg := db.CreateReservationParams{}
	err = arg.Unmarshal(reservationData)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// generate reservation code to add to parameters
	arg.Code, err = util.GenerateReservationCode(reservation.LastName, ReservationCodeLenght)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// insert new reservation into database
	result, err := s.DatabaseStore.CreateReservation(ctx, arg)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// update reservation with database result
	err = reservation.Unmarshal(result.Marshal())
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// load reservation into session data
	app.Session.Put(r.Context(), "reservation", reservation)

	// redirecting to summery page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (s *Server) ReservationSummary(w http.ResponseWriter, r *http.Request) {
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
