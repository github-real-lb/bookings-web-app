package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/github-real-lb/bookings-web-app/util/forms"
	"github.com/go-chi/chi/v5"
)

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
		app.LogServerError(w, err)
	}
}

// RoomsHandler is the GET "/rooms/{index}" page handler
func (s *Server) RoomsHandler(w http.ResponseWriter, r *http.Request) {
	// if no id paramater exists in URL render a new page
	if chi.URLParam(r, "index") == "list" {
		//TODO: change this to UI input
		rooms, err := s.ListRooms(10, 0)
		if err != nil {
			app.LogServerError(w, err)
			return
		}

		app.Session.Put(r.Context(), "rooms", rooms)

		err = RenderTemplate(w, r, "rooms.page.gohtml", &TemplateData{
			Data: map[string]any{
				"rooms": rooms,
			},
		})
		if err != nil {
			app.LogServerError(w, err)
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
		app.LogServerError(w, err)
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
		app.LogServerError(w, err)
	}
}

// PostSearchRoomAvailabilityHandler is the POST "/search-room-availability" page handler
// It is fetched by the room.page and excpect a json response
func (s *Server) PostSearchRoomAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	room, ok := app.Session.Get(r.Context(), "room").(Room)
	if !ok {
		http.Redirect(w, r, "/rooms/list", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// define the json response
	var resp struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	if ok = form.Required("start_date"); !ok {
		resp.OK = false
		resp.Message = form.Errors.Get("start_date")
	} else if ok = form.Required("end_date"); !ok {
		resp.OK = false
		resp.Message = form.Errors.Get("end_date")
	} else if ok = form.CheckDateRange("start_date", "end_date"); !ok {
		resp.OK = false
		resp.Message = form.Errors.Get("start_date")
	}

	// returns response if form data are invalid
	if !form.Valid() {
		bs, err := json.Marshal(resp)
		if err != nil {
			app.LogServerError(w, err)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			app.LogError(err)
			return
		}
		return
	}

	// parse form's data to reservation
	var reservation Reservation
	reservation.Room = room
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// check if room is available
	ok, err = s.CheckRoomAvailability(reservation)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// return response if room is not available
	if !ok {
		resp.OK = false
		resp.Message = "Room is unavailable. PLease try different dates."

		bs, err := json.Marshal(resp)
		if err != nil {
			app.LogServerError(w, err)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			app.LogError(err)
			return
		}
		return
	}

	resp.OK = true

	bs, err := json.Marshal(resp)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// load reservation to session data
	app.Session.Put(r.Context(), "reservation", reservation)

	_, err = w.Write(bs)
	if err != nil {
		app.LogError(err)
	}
}

// ContactHandler is the GET "/contact" page handler
func (s *Server) ContactHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "contact.page.gohtml", &TemplateData{})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// AvailabilityHandler is the GET "/available-rooms-search" page handler
func (s *Server) AvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "available-rooms-search.page.gohtml", &TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		app.LogServerError(w, err)
	}
}

// PostAvailability is the POST "/available-rooms-search" page handler
func (s *Server) PostAvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
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

	// get list of availabe rooms
	rooms, err := s.ListAvailableRooms(reservation)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// check if there are rooms availabe
	if len(rooms) == 0 {
		app.Session.Put(r.Context(), "warning", "No rooms are availabe. Please try different dates.")
		err = RenderTemplate(w, r, "available-rooms-search.page.gohtml", &TemplateData{
			Form: form,
		})
		if err != nil {
			app.LogServerError(w, err)
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
		app.LogError(errors.New("cannot get available rooms list from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
			app.LogServerError(w, err)
		}
		return
	}

	// get room id from URL
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// get reservation data from session
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		app.LogError(errors.New("cannot get reservation from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		app.LogError(errors.New("cannot get reservation data from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := RenderTemplate(w, r, "make-reservation.page.gohtml", &TemplateData{
		StringMap: map[string]string{
			"start_date": reservation.StartDate.Format("2006-01-02"),
			"end_date":   reservation.EndDate.Format("2006-01-02"),
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
		app.LogError(errors.New("cannot get reservation data from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.LogServerError(w, err)
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
				"start_date": reservation.StartDate.Format("2006-01-02"),
				"end_date":   reservation.StartDate.Format("2006-01-02"),
			},
			Data: map[string]any{
				"reservation": reservation,
			},
			Form: form,
		})
		if err != nil {
			app.LogServerError(w, err)
		}
		return
	}

	// parse form's data to reservation
	err = reservation.Unmarshal(form.Marshal())
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// generate reservation code
	reservation.Code, err = util.GenerateReservationCode(reservation.LastName, ReservationCodeLenght)
	if err != nil {
		app.LogServerError(w, err)
		return
	}

	// insert reservation into database
	// TODO: the 1 should be replaced with database ENUM input of restriction id
	err = s.CreateReservation(&reservation, 1)
	if err != nil {
		app.LogServerError(w, err)
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
		app.LogError(errors.New("cannot get reservation data from the session"))
		app.Session.Put(r.Context(), "error", "No reservation exists. Please make a reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove reservation and rooms data from session
	app.Session.Remove(r.Context(), "reservation")
	app.Session.Remove(r.Context(), "rooms")

	err := RenderTemplate(w, r, "reservation-summary.page.gohtml", &TemplateData{
		StringMap: map[string]string{
			"start_date": reservation.StartDate.Format("2006-01-02"),
			"end_date":   reservation.EndDate.Format("2006-01-02"),
		},
		Data: map[string]any{
			"reservation": reservation,
		},
	})

	if err != nil {
		app.LogServerError(w, err)
	}
}
