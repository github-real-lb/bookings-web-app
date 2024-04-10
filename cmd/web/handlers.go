package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/forms"
	"github.com/go-chi/chi/v5"
)

// LimitRoomsPerPage sets the maximum number of rooms to display on a page
const LimitRoomsPerPage = 10

// HomeHandler is the GET "/" home page handler
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "home.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// AboutHandler is the GET "/about" page handler
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "about.page.gohtml", &TemplateData{})
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to render about.page.gohtml template", err, "/")
	}
}

// RoomsHandler is the GET "/rooms/{index}" page handler
func (s *Server) RoomsHandler(w http.ResponseWriter, r *http.Request) {
	// if no id paramater exists in URL render a new page
	if chi.URLParam(r, "index") == "list" {
		//TODO: change the offset to request input
		rooms, err := s.ListRooms(LimitRoomsPerPage, 0)
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to load rooms from database", err, "/")
			return
		}

		app.Session.Put(r.Context(), "rooms", rooms)

		err = RenderTemplate(w, r, "rooms.page.gohtml", &TemplateData{
			Data: map[string]any{
				"rooms": rooms,
			},
		})
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to render rooms.page.gohtml template", err, "/")
		}
		return
	}

	rooms, ok := app.Session.Get(r.Context(), "rooms").(Rooms)
	if !ok {
		http.Redirect(w, r, "/rooms/list", http.StatusTemporaryRedirect)
		return
	}

	// get room id from URL
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		http.Redirect(w, r, "/rooms/list", http.StatusTemporaryRedirect)

		return
	}

	// check if index is out of scope
	if index >= len(rooms) {
		http.Redirect(w, r, "/rooms/list", http.StatusTemporaryRedirect)
		return
	}

	// put selected room data to session
	room := rooms[index]
	app.Session.Put(r.Context(), "room", room)

	// remove rooms data from session
	app.Session.Remove(r.Context(), "rooms")

	// create redirect url
	url := strings.ReplaceAll(room.Name, "'", "")
	url = strings.ReplaceAll(url, " ", "-")
	url = fmt.Sprint("/rooms/room/", url)

	//redirecting to make-reservation page
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// RoomHandler is the GET "/rooms/room/{name}" page handler
func (s *Server) RoomHandler(w http.ResponseWriter, r *http.Request) {
	room, ok := app.Session.Get(r.Context(), "room").(Room)
	if !ok {
		http.Redirect(w, r, "/rooms/list", http.StatusTemporaryRedirect)
		return
	}

	err := RenderTemplate(w, r, "room.page.gohtml", &TemplateData{
		Data: map[string]any{
			"room": room,
		},
	})
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to render rooms.page.gohtml template", err, "/rooms/list")
	}
}

// define the type of json response
type SearchRoomAvailabilityResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// PostSearchRoomAvailabilityHandler is the POST "/search-room-availability" page handler
// It is fetched by the room.page and excpect a json response
func (s *Server) PostSearchRoomAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	room, ok := app.Session.Get(r.Context(), "room").(Room)
	if !ok {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})
		s.LogError(r, "unable to get room from session", errors.New("ERROR: wrong routing"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})
		s.LogError(r, "unable to parse form", err)
		return
	}

	// create a new form with data and validate the form
	var errMsg string
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	if ok = form.Required("start_date"); !ok {
		errMsg = form.Errors.Get("start_date")
	} else if ok = form.Required("end_date"); !ok {
		errMsg = form.Errors.Get("end_date")
	} else if ok = form.CheckDateRange("start_date", "end_date"); !ok {
		errMsg = form.Errors.Get("end_date")
	}

	// returns response if form data are invalid
	if !form.Valid() {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:      false,
			Message: errMsg,
		})
		return
	}

	// parse form's data to reservation
	var reservation Reservation
	reservation.Room = room
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})
		s.LogError(r, "unable to unmarshal reservation", err)
		return
	}

	// check if room is available
	ok, err = s.CheckRoomAvailability(reservation)
	if err != nil {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})
		s.LogError(r, "unable to check room availability", err)
		return
	}

	if ok {
		// load reservation to session data
		app.Session.Put(r.Context(), "reservation", reservation)

		// write the json response
		if err = s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{OK: true}); err != nil {
			return
		}
	} else {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:      false,
			Message: "Room is unavailable. PLease try different dates.",
		})
	}
}

// ContactHandler is the GET "/contact" page handler
func (s *Server) ContactHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "contact.page.gohtml", &TemplateData{})
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to render contact.page.gohtml template", err, "/")
	}
}

// AvailabilityHandler is the GET "/available-rooms-search" page handler
func (s *Server) AvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "available-rooms-search.page.gohtml", &TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to render available-rooms-search.page.gohtml template", err, "/")
	}
}

// PostAvailability is the POST "/available-rooms-search" page handler
func (s *Server) PostAvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to parse form", err, "/available-rooms-search")
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	form.Required("start_date", "end_date")
	form.CheckDateRange("start_date", "end_date")

	if !form.Valid() {
		err = RenderTemplate(w, r, "available-rooms-search.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to render available-rooms-search.page.gohtml template", err, "/")
		}
		return
	}

	// parse form's data to reservation
	var reservation Reservation
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to parse form data for reservation ", err, "/available-rooms-search")
		return
	}

	// get list of available rooms
	rooms, err := s.ListAvailableRooms(reservation, LimitRoomsPerPage, 0)
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to load available rooms", err, "/available-rooms-search")
		return
	}

	// check if there are rooms available
	if len(rooms) == 0 {
		app.Session.Put(r.Context(), "warning", "No rooms are available. Please try different dates.")
		err = RenderTemplate(w, r, "available-rooms-search.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to render available-rooms-search.page.gohtml template", err, "/")
		}
		return
	}

	// load reservation to session data
	app.Session.Put(r.Context(), "reservation", reservation)
	app.Session.Put(r.Context(), "rooms", rooms)

	// redirecting to choose-room page
	http.Redirect(w, r, "/available-rooms/available", http.StatusSeeOther)
}

// ChooseRoomHandler is the GET "/available-rooms/{index}" page handler
func (s *Server) AvailableRoomsListHandler(w http.ResponseWriter, r *http.Request) {
	// get available rooms data from session
	rooms, ok := app.Session.Get(r.Context(), "rooms").(Rooms)
	if !ok {
		s.LogErrorAndRedirect(w, r, "No reservation exists. Please make a reservation.", errors.New("ERROR: wrong routing"), "/available-rooms-search")
		return
	}

	// if no id paramater exists in URL render a new page
	if chi.URLParam(r, "index") == "available" {
		err := RenderTemplate(w, r, "available-rooms.page.gohtml", &TemplateData{
			Data: map[string]any{
				"rooms": rooms,
			},
		})
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to render available-rooms.page.gohtml template", err, "/")
		}
		return
	}

	// get room id from URL
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to convert index parameter to integer", err, "/available-rooms/available")
		return
	}

	// get reservation data from session
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		s.LogErrorAndRedirect(w, r, "No reservation exists. Please make a reservation.", errors.New("ERROR: wrong routing"), "/")
		return
	}

	reservation.Room.ID = rooms[index].ID
	reservation.Room.Name = rooms[index].Name
	reservation.Room.Description = rooms[index].Description
	reservation.Room.ImageFilename = rooms[index].ImageFilename
	app.Session.Put(r.Context(), "reservation", reservation)

	// redirecting to make-reservation page
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ReservationHandler is the GET "/make-reservation" page handler
func (s *Server) MakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		s.LogErrorAndRedirect(w, r, "No reservation exists. Please make a reservation.", errors.New("ERROR: wrong routing"), "/")
		return
	}

	err := RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
		StringMap: map[string]string{
			"start_date": reservation.StartDate.Format(config.DateLayout),
			"end_date":   reservation.EndDate.Format(config.DateLayout),
		},
		Data: map[string]any{
			"reservation": reservation,
		},
		Form: forms.New(nil),
	})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostReservationHandler is the POST "/make-reservation" page handler
func (s *Server) PostMakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		s.LogErrorAndRedirect(w, r, "No reservation exists. Please make a reservation.", errors.New("ERROR: wrong routing"), "/")
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to parse form", err, "/make-reservation")
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	form.Required("first_name", "last_name", "email")
	form.CheckMinLenght("first_name", 3)
	form.CheckMinLenght("last_name", 3)
	form.CheckEmail("email")

	if !form.Valid() {
		err = RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
			StringMap: map[string]string{
				"start_date": reservation.StartDate.Format(config.DateLayout),
				"end_date":   reservation.StartDate.Format(config.DateLayout),
			},
			Data: map[string]any{
				"reservation": reservation,
			},
			Form: form,
		})
		if err != nil {
			s.LogErrorAndRedirect(w, r, "unable to render make-reservation.page.gohtml", err, "/make-reservation")
		}
		return
	}

	// parse form's data to reservation
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to parse form data into reservation", err, "/make-reservation")
		return
	}

	// generate reservation code
	err = reservation.GenerateReservationCode()
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to generate reservation code", err, "/make-reservation")
		return
	}

	// insert reservation into database
	err = s.CreateReservation(&reservation)
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to create reservation", err, "/make-reservation")
		return
	}

	// load reservation data into session
	app.Session.Put(r.Context(), "reservation", reservation)

	// redirecting to summery page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummaryHandler is the GET "/reservation-summery" page handler
func (s *Server) ReservationSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// get reservation data from session
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		s.LogErrorAndRedirect(w, r, "No reservation exists. Please make a reservation.", errors.New("ERROR: wrong routing"), "/")
		return
	}

	// remove reservation and rooms data from session
	app.Session.Remove(r.Context(), "reservation")
	app.Session.Remove(r.Context(), "rooms")

	err := RenderTemplate(w, r, "reservation-summary.page.gohtml", &TemplateData{
		StringMap: map[string]string{
			"start_date": reservation.StartDate.Format(config.DateLayout),
			"end_date":   reservation.EndDate.Format(config.DateLayout),
		},
		Data: map[string]any{
			"reservation": reservation,
		},
	})
	if err != nil {
		s.LogErrorAndRedirect(w, r, "unable to create reservation", err, "/")
	}
}
