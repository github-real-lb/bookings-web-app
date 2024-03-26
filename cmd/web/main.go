package main

import (
	"context"
	"fmt"
	"log"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/webapp"
)

const ReservationCodeLenght = 6

func main() {
	// initializing application
	err := InitializeApp(webapp.TestingMode)
	if err != nil {
		log.Fatal("Error initializing application:", err)
	}

	// connecting to postgres database
	store, err := db.NewPostgresDBStore(app.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer store.(*db.PostgresDBStore).Pool.Close()

	// start the server
	server := NewServer(store)

	AddDataToDB(server)

	fmt.Println("Starting web server on,", app.ServerAddress)
	err = server.Router.ListenAndServe()
	log.Fatal(err)
}

func AddDataToDB(s *Server) {
	s.DatabaseStore.CreateRoom(context.Background(), "General's Quarters")
	s.DatabaseStore.CreateRoom(context.Background(), "Major's Suite")
}
