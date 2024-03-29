package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/github-real-lb/bookings-web-app/db"
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
func (s *Server) SearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostAvailability is the search-availability page post handler
func (s *Server) PostSearchAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)

	// validate form
	form.TrimSpaces()
	form.Required("start_date", "end_date")
	form.CheckDateRange("start_date", "end_date")

	if !form.Valid() {
		err = RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			app.LogServerError(w, err)
		}
		return
	}

	// parse form's data to query arguments
	var arg db.ListAvailableRoomsParams
	arg.Unmarshal(form.Marshal())
	rooms, err := s.DatabaseStore.ListAvailableRooms(context.Background(), arg)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		app.Session.Put(r.Context(), "warning", "No rooms are availabe. Please try different dates.")
		err = RenderTemplate(w, r, "search-availability.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			app.LogServerError(w, err)
		}
		return
	}

	// parse form's data to reservation
	var reservation Reservation
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		app.LogServerError(w, err)
	}

	app.Session.Put(r.Context(), "reservation", reservation)
	err = RenderTemplate(w, r, "choose-room.page.gohtml", &TemplateData{
		Data: map[string]any{
			"rooms": rooms,
		},
	})
	if err != nil {
		app.LogServerError(w, err)
	}

}

// ReservationHandler is the make-reservation page handler
func (s *Server) MakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostReservationHandler is the make-reservation post page handler
func (s *Server) PostMakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)

	//TODO: replace this code part with code loading the data from the session
	form.Add("room_id", "1")

	//TODO: update validation to include all fields
	form.TrimSpaces()
	form.Required("first_name", "last_name", "email")
	form.CheckMinLenght("first_name", 3)
	form.CheckMinLenght("last_name", 3)
	form.CheckEmail("email")

	if !form.Valid() {
		err = RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			app.LogServerError(w, err)
		}
		return
	}

	// parse form's data to reservation
	var reservation Reservation
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		app.LogServerError(w, err)
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
