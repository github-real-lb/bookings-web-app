package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReservationConfirmationMail(t *testing.T) {
	r := randomReservation()
	mailData, err := CreateReservationConfirmationMail(r)
	require.NoError(t, err)
	assert.Equal(t, r.Email, mailData.To)
	assert.Equal(t, app.Listing.Email, mailData.From)
	assert.Equal(t, fmt.Sprintf("Confirmation Notice for Reservation %s", r.Code), mailData.Subject)
	assert.NotEmpty(t, mailData.Content)
}
