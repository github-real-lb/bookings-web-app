package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/github-real-lb/bookings-web-app/db"
	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/mailer"
)

func main() {
	// initialize application
	err := InitializeApp(config.DevelopmentMode)
	if err != nil {
		log.Fatal("Error initializing application:", err)
	}

	// connect to postgres database
	dbStore, err := db.NewPostgresDBStore(app.DBConnectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer dbStore.(*db.PostgresDBStore).DBConnPool.Close()

	// start listenning for errors
	app.Logger.ListenAndLogErrors()
	defer close(app.Logger.ErrorChannel)

	// create email channel and start listening for data
	app.MailerChan = mailer.GetMailerChannel()
	defer close(app.MailerChan)
	mailer.Listen(app.Logger.ErrorChannel)

	// create a new server and starting it in a separate goroutine
	server := NewServer(dbStore)
	go server.Start()

	// Listen for interrupt signal (Ctrl+C) or SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)    // Ctrl+C
	signal.Notify(stop, syscall.SIGTERM) // SIGTERM

	// block until a stop signal is received
	<-stop

	// stop logger and server
	app.Logger.Shutdown()
	server.Stop()
}
