package main

import (
	"fmt"

	"github.com/github-real-lb/bookings-web-app/util/config"
	"github.com/github-real-lb/bookings-web-app/util/mailers"
)

// CreateReservationNotificationMail creates reservation confirmation mail
func CreateReservationConfirmationMail(r Reservation) (mailers.MailData, error) {
	var err error

	// create reservation confirmation email
	data := mailers.MailData{
		To:      r.Email,
		From:    app.Listing.Email,
		Subject: fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code),
	}

	data.Content, err = RenderMailTemplate("reservation-confirmation.mail.gohtml", &TemplateData{
		Data: map[string]any{
			"start_date":  r.StartDate.Format(config.DateLayout),
			"end_date":    r.EndDate.Format(config.DateLayout),
			"reservation": r,
		},
	})

	return data, err
}
