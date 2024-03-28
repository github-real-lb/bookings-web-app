package main

import (
	"fmt"
	"log"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/config"
)

const ReservationCodeLenght = 6

func main() {
	// initializing application
	err := InitializeApp(config.DevelopmentMode)
	if err != nil {
		log.Fatal("Error initializing application:", err)
	}

	// connecting to postgres database
	store, err := db.NewPostgresDBStore(app.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer store.(*db.PostgresDBStore).DBConnPool.Close()

	// start the server
	server := NewServer(store)

	fmt.Println("Starting web server on,", app.ServerAddress)
	err = server.Router.ListenAndServe()
	log.Fatal(err)
}
