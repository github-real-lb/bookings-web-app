package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/util/forms"
)

// HomeHandler is the home page handler
func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "home.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// AboutHandler is the about page handler
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "about.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// ReservationHandler is the generals-quarters room page handler
func (s *Server) GeneralsHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "generals.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// MajorsHandler is the majors-suite room page handler
func (s *Server) MajorsHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "majors.room.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// ContactHandler is the contact page handler
func (s *Server) ContactHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "contact.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// AvailabilityHandler is the search-availability page handler
func (s *Server) AvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostAvailability is the search-availability page post handler
func (s *Server) PostAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
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

// PostAvailabilityJsonHandler handles requests for availability and sends JSON response
func (s *Server) PostAvailabilityJsonHandler(w http.ResponseWriter, r *http.Request) {
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

// ReservationHandler is the make-reservation page handler
func (s *Server) ReservationHandler(w http.ResponseWriter, r *http.Request) {
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

// PostReservationHandler is the make-reservation post page handler
func (s *Server) PostReservationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	var form *forms.Form
	var reservation Reservation

	// create a new form with data and validate the form
	form = forms.New(r.PostForm)

	//TODO: replace this code part with code loading the data from the session
	form.Add("room_id", "1")

	//TODO: update validation to include all fields
	form.TrimSpaces()
	form.Required("first_name", "last_name", "email")
	form.MinLenght("first_name", 3)
	form.MinLenght("last_name", 3)
	form.IsEmailValid("email")

	// parse form's data to reservation
	err = reservation.Unmarshal(form.Marshal())
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

	// insert reservation into database
	// TODO: the 1 should be replaced with UI input of restriction id
	err = s.CreateReservation(&reservation, 1)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// load reservation into session data
	app.Session.Put(r.Context(), "reservation", reservation)

	// redirecting to summery page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (s *Server) ReservationSummaryHandler(w http.ResponseWriter, r *http.Request) {
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
