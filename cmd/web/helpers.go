package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
)

const ReservationCodeLenght = 7

// GenerateReservationCode generate the reservation code.
func (r *Reservation) GenerateReservationCode() {
	// concatenating the current time with the reservation last name
	s := fmt.Sprintf("%v %s %v %v", util.RandomDatetime().Format(time.RFC3339Nano), r.LastName, r.StartDate, r.EndDate)

	// Generate SHA-256 hash of the concatenated string
	hash := sha256.New()
	hash.Write([]byte(s))

	// Generate the SHA256 checksum of the data written so far and convert to hexadecimal string
	hashString := hex.EncodeToString(hash.Sum(nil))

	// build code string
	code := make([]byte, ReservationCodeLenght)
	digitsFound := 0
	digitsMax := ReservationCodeLenght / 2
	lettersFound := 0
	lettersMax := ReservationCodeLenght - digitsMax

	for _, v := range []byte(hashString) {
		if (digitsFound < digitsMax) && (v >= 49 && v <= 57) {
			// adds digits to code if not enought digits were found and if v is a digit between 1-9
			code[digitsFound*2+1] = v
			digitsFound++
		} else if (lettersFound < lettersMax) && ((v >= 97 && v <= 104) || (v >= 106 && v <= 110) || (v >= 112 && v <= 122)) {
			// adds letters to code if not enought letters were found check if v is a letter except for 'i' or 'o'
			code[lettersFound*2] = v - 32
			lettersFound++
		}

		if digitsFound+lettersFound == ReservationCodeLenght {
			r.Code = string(code)
			return
		}
	}
}

func IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}
