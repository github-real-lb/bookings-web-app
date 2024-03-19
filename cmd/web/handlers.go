package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/internal/forms"
	"github.com/github-real-lb/bookings-web-app/internal/models"
	"github.com/github-real-lb/bookings-web-app/internal/render"
)

// Home is the home page handler
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (server *Server) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Reservation is the generals-quarters room page handler
func (server *Server) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.room.page.gohtml", &models.TemplateData{})
}

// Majors is the majors-suite room page handler
func (server *Server) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.room.page.gohtml", &models.TemplateData{})
}

// Contact is the contact page handler
func (server *Server) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// Availability is the search-availability page handler
func (server *Server) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability is the search-availability page post handler
func (server *Server) PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start_date")
	endDate := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", startDate, endDate)))
}

type jsonResponse struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// AvailabilityJSON handles requests for availability and sends JSON response
func (server *Server) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	startDate, err := time.Parse("2006-01-02", r.Form.Get("start_date"))
	if err != nil {
		log.Println(err)
	}
	endDate, err := time.Parse("2006-01-02", r.Form.Get("end_date"))
	if err != nil {
		log.Println(err)
	}

	resp := jsonResponse{
		StartDate: startDate,
		EndDate:   endDate,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation is the make-reservation page handler
func (server *Server) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: map[string]any{
			"reservation": models.Reservation{},
		},
	})
}

// PostReservation is the make-reservation post page handler
func (server *Server) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLenght("first_name", 3)
	form.MinLenght("last_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: map[string]any{
				"reservation": reservation,
			},
		})

		return
	}

	app.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (server *Server) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get reservation from the session.")
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	app.Session.Remove(r.Context(), "reservation")

	render.RenderTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: map[string]any{
			"reservation": reservation,
		},
	})
}
