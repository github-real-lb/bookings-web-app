package main

import (
	"fmt"
	"log"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/config"
)

func main() {
	// initializing application
	err := InitializeApp(config.DevelopmentMode)
	if err != nil {
		log.Fatal("Error initializing application:", err)
	}

	// connecting to postgres database
	dbStore, err := db.NewPostgresDBStore(app.DBConnectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer dbStore.(*db.PostgresDBStore).DBConnPool.Close()

	// start the server
	server := NewServer(dbStore)

	fmt.Println("Starting web server on,", app.ServerAddress)
	err = server.Router.ListenAndServe()
	log.Fatal(err)
}
