package main

import (
	"errors"
	"fmt"
	"log"
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
	err := s.Renderer.RenderGoHtmlPageTemplate(w, r, "home.page.gohtml", &TemplateData{})
	if err != nil {
		sErr := CreateServerError(ErrorRenderTemplate, r.URL.Path, err)
		s.LogError(sErr)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// AboutHandler is the GET "/about" page handler
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request) {
	s.Render(w, r, "about.page.gohtml", &TemplateData{}, "/")
}

// RoomsHandler is the GET "/rooms/{index}" page handler
func (s *Server) RoomsHandler(w http.ResponseWriter, r *http.Request) {
	// if no id paramater exists in URL render a new page
	if chi.URLParam(r, "index") == "list" {
		//TODO: change the offset to request input
		rooms, err := s.ListRooms(LimitRoomsPerPage, 0)
		if err != nil {
			sErr := ServerError{
				Prompt: "Unable to load rooms from database.",
				URL:    r.URL.Path,
				Err:    err,
			}
			s.LogErrorAndRedirect(w, r, sErr, "/")
			return
		}

		app.Session.Put(r.Context(), "rooms", rooms)

		s.Render(w, r, "rooms.page.gohtml",
			&TemplateData{
				Data: map[string]any{"rooms": rooms},
			}, "/")
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

	s.Render(w, r, "room.page.gohtml",
		&TemplateData{
			Data: map[string]any{"room": room},
		}, "/rooms/list")
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

		s.LogError(ServerError{
			Prompt: "Unable to get room from session.",
			URL:    r.URL.Path,
			Err:    errors.New("wrong routing"),
		})
		return
	}

	err := r.ParseForm()
	if err != nil {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})

		sErr := CreateServerError(ErrorParseForm, r.URL.Path, err)
		s.LogError(sErr)
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
	rsv := Reservation{}
	form.GetValue("start_date", &rsv.StartDate)
	form.GetValue("end_date", &rsv.EndDate)
	rsv.RoomID = room.ID
	rsv.Room = room

	// check if room is available
	ok, err = s.CheckRoomAvailability(rsv.RoomID, rsv.StartDate, rsv.EndDate)
	if err != nil {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:    false,
			Error: "Internal Error. Please reload and try again.",
		})

		s.LogError(ServerError{
			Prompt: "Unable to check room availability.",
			URL:    r.URL.Path,
			Err:    err,
		})
		return
	}

	if ok {
		// load reservation to session data
		app.Session.Put(r.Context(), "reservation", rsv)

		// write the json response
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{OK: true})
	} else {
		s.ResponseJSON(w, r, SearchRoomAvailabilityResponse{
			OK:      false,
			Message: "Room is unavailable. PLease try different dates.",
		})
	}
}

// ContactHandler is the GET "/contact" page handler
func (s *Server) ContactHandler(w http.ResponseWriter, r *http.Request) {
	s.Render(w, r, "contact.page.gohtml", &TemplateData{}, "/")
}

// AvailabilityHandler is the GET "/available-rooms-search" page handler
func (s *Server) AvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	s.Render(w, r, "available-rooms-search.page.gohtml",
		&TemplateData{Form: forms.New(nil)}, "/")
}

// PostAvailability is the POST "/available-rooms-search" page handler
func (s *Server) PostAvailableRoomsSearchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		sErr := CreateServerError(ErrorParseForm, r.URL.Path, err)
		s.LogErrorAndRedirect(w, r, sErr, "/available-rooms-search")
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	form.Required("start_date", "end_date")
	form.CheckDateRange("start_date", "end_date")

	if !form.Valid() {
		s.Render(w, r, "available-rooms-search.page.gohtml",
			&TemplateData{Form: form}, "/")
		return
	}

	// parse form's data to reservation
	rsv := Reservation{}
	form.GetValue("start_date", &rsv.StartDate)
	form.GetValue("end_date", &rsv.EndDate)

	// get list of available rooms
	rooms, err := s.ListAvailableRooms(LimitRoomsPerPage, 0, rsv.StartDate, rsv.EndDate)
	if err != nil {
		sErr := ServerError{
			Prompt: "Unable to load available rooms.",
			URL:    r.URL.Path,
			Err:    err,
		}
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	// check if there are rooms available
	if len(rooms) == 0 {
		app.Session.Put(r.Context(), "warning", "No rooms are available. Please try different dates.")
		s.Render(w, r, "available-rooms-search.page.gohtml",
			&TemplateData{Form: form}, "/")
		return
	}

	// load reservation to session data
	app.Session.Put(r.Context(), "reservation", rsv)
	app.Session.Put(r.Context(), "rooms", rooms)

	// redirecting to choose-room page
	http.Redirect(w, r, "/available-rooms/available", http.StatusSeeOther)
}

// ChooseRoomHandler is the GET "/available-rooms/{index}" page handler
func (s *Server) AvailableRoomsListHandler(w http.ResponseWriter, r *http.Request) {
	// get available rooms data from session
	rooms, ok := app.Session.Get(r.Context(), "rooms").(Rooms)
	if !ok {
		sErr := CreateServerError(ErrorMissingReservation, r.URL.Path, nil)
		s.LogErrorAndRedirect(w, r, sErr, "/available-rooms-search")
		return
	}

	// if no id paramater exists in URL render a new page
	if chi.URLParam(r, "index") == "available" {
		s.Render(w, r, "available-rooms.page.gohtml",
			&TemplateData{
				Data: map[string]any{"rooms": rooms},
			}, "/")
		return
	}

	// get room id from URL
	index, err := strconv.Atoi(chi.URLParam(r, "index"))
	if err != nil {
		sErr := ServerError{
			Prompt: "Unable to convert index parameter to integer.",
			URL:    r.URL.Path,
			Err:    err,
		}
		s.LogErrorAndRedirect(w, r, sErr, "/available-rooms/available")
		return
	}

	// get reservation data from session
	rsv, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		sErr := CreateServerError(ErrorMissingReservation, r.URL.Path, nil)
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	rsv.RoomID = rooms[index].ID
	rsv.Room = rooms[index]
	app.Session.Put(r.Context(), "reservation", rsv)

	// redirecting to make-reservation page
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ReservationHandler is the GET "/make-reservation" page handler
func (s *Server) MakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		sErr := CreateServerError(ErrorMissingReservation, r.URL.Path, nil)
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	s.Render(w, r, "make-reservation.page.gohtml",
		&TemplateData{
			Data: map[string]any{
				"start_date":  reservation.StartDate.Format(config.DateLayout),
				"end_date":    reservation.EndDate.Format(config.DateLayout),
				"reservation": reservation,
			},
			Form: forms.New(nil),
		}, "/")
}

// PostReservationHandler is the POST "/make-reservation" page handler
func (s *Server) PostMakeReservationHandler(w http.ResponseWriter, r *http.Request) {
	rsv, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		sErr := CreateServerError(ErrorMissingReservation, r.URL.Path, nil)
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	err := r.ParseForm()
	if err != nil {
		sErr := CreateServerError(ErrorParseForm, r.URL.Path, err)
		s.LogErrorAndRedirect(w, r, sErr, "/make-reservation")
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	form.Required("first_name", "last_name", "email")
	form.CheckMinLenght("first_name", 3)
	form.CheckMinLenght("last_name", 3)
	form.CheckEmail("email")

	log.Println("TODO: validate phone and notes even if not required")

	if !form.Valid() {
		s.Render(w, r, "make-reservation.page.gohtml",
			&TemplateData{
				Data: map[string]any{
					"start_date":  rsv.StartDate.Format(config.DateLayout),
					"end_date":    rsv.StartDate.Format(config.DateLayout),
					"reservation": rsv,
				},
				Form: form,
			}, "/make-reservation")
		return
	}

	// parse form's data to reservation
	form.GetValue("first_name", &rsv.FirstName)
	form.GetValue("last_name", &rsv.LastName)
	form.GetValue("email", &rsv.Email)
	form.GetValue("phone", &rsv.Phone)
	form.GetValue("notes", &rsv.Notes)

	// generate reservation code
	rsv.GenerateReservationCode()

	// insert reservation into database
	err = s.CreateReservation(rsv)
	if err != nil {
		sErr := ServerError{
			Prompt: "Unable to create reservation.",
			URL:    r.URL.Path,
			Err:    err,
		}
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	// load reservation data into session
	app.Session.Put(r.Context(), "reservation", rsv)

	data, err := s.Renderer.CreateReservationConfirmationMail(rsv)
	if err != nil {
		sErr := ServerError{
			Prompt: "Unable to render confirmation email.",
			URL:    r.URL.Path,
			Err:    err,
		}
		s.LogErrorAndRedirect(w, r, sErr, "/reservation-summary")
		return
	}

	// send reservation notification email to guest and log
	s.SendMail(data)
	s.LogInfo(fmt.Sprintf("MAIL confirmation notice sent to %s", data.To))

	// send reservation notification email to admin and log
	data.To = "admin@listingdomain.com" // TODO: this should be update based on app admin setting
	s.SendMail(data)
	s.LogInfo(fmt.Sprintf("MAIL confirmation notice sent to %s", data.To))

	// redirecting to summery page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummaryHandler is the GET "/reservation-summery" page handler
func (s *Server) ReservationSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// get reservation data from session
	reservation, ok := app.Session.Get(r.Context(), "reservation").(Reservation)
	if !ok {
		sErr := CreateServerError(ErrorMissingReservation, r.URL.Path, nil)
		s.LogErrorAndRedirect(w, r, sErr, "/")
		return
	}

	// remove reservation and rooms data from session
	app.Session.Remove(r.Context(), "reservation")
	app.Session.Remove(r.Context(), "rooms")

	s.Render(w, r, "reservation-summary.page.gohtml",
		&TemplateData{
			Data: map[string]any{
				"start_date":  reservation.StartDate.Format(config.DateLayout),
				"end_date":    reservation.EndDate.Format(config.DateLayout),
				"reservation": reservation,
			},
		}, "/")
}

// LoginHandler is the GET "/user/login" page handler
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	s.Render(w, r, "login.page.gohtml",
		&TemplateData{Form: forms.New(nil)}, "/")
}

// PostLoginHandler is the POST "/user/login" page handler
func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		sErr := CreateServerError(ErrorParseForm, r.URL.Path, err)
		s.LogErrorAndRedirect(w, r, sErr, "/user/login")
		return
	}

	// create a new form with data and validate the form
	form := forms.New(r.PostForm)
	form.TrimSpaces()
	form.Required("email", "password")
	form.CheckEmail("email")
	form.CheckPassword("password")

	if !form.Valid() {
		s.Render(w, r, "login.page.gohtml",
			&TemplateData{Form: form}, "/")
		return
	}

	// authenticate the user
	userID, err := s.AuthenticateUser(form.Get("email"), form.Get("password"))
	if err != nil {
		form.Del("password")

		s.LogInfo(fmt.Sprintf("Unsuccessful login with email %s", form.Get("email")))

		s.Render(w, r, "login.page.gohtml",
			&TemplateData{
				Form:  forms.New(nil),
				Error: "Invalid Email and/or Password. Please try again.",
			}, "/")

		return
	}

	app.Session.Put(r.Context(), "user_id", userID)
	app.Session.Put(r.Context(), "flash", "Successfully logged in!")
	s.LogInfo(fmt.Sprintf("Successful login by user %d", userID))

	// redirecting to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler is the GET "/user/logout" page handler
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())
	app.Session.Put(r.Context(), "flash", "Successfully logged out!")

	// redirecting to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboardHandler is the GET "/user/admin" page handler
func (s *Server) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	s.Render(w, r, "admin-dashboard.page.gohtml", &TemplateData{}, "/")
}
