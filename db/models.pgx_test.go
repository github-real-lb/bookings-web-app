package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/require"
)

func TestReservation_MarshalAndUnmarhsal(t *testing.T) {
	data := make(map[string]string)
	data["id"] = fmt.Sprint(util.RandomID())
	data["code"] = util.RandomString(ReservationCodeLenght)
	data["first_name"] = util.RandomName()
	data["last_name"] = util.RandomName()
	data["email"] = util.RandomEmail()
	data["phone"] = util.RandomPhone()
	data["start_date"] = util.RandomDate().Format("2006-01-02")
	data["end_date"] = util.RandomDate().Format("2006-01-02")
	data["room_id"] = fmt.Sprint(util.RandomID())
	data["notes"] = util.RandomNote()
	data["created_at"] = util.RandomDatetime().Format(time.RFC3339)
	data["updated_at"] = util.RandomDatetime().Format(time.RFC3339)

	r := Reservation{}
	err := r.Unmarshal(data)
	require.NoError(t, err)
	require.Equal(t, data["id"], fmt.Sprint(r.ID))

	require.Equal(t, data["code"], r.Code)
	require.Equal(t, data["first_name"], r.FirstName)
	require.Equal(t, data["last_name"], r.LastName)
	require.Equal(t, data["email"], r.Email)

	require.Equal(t, data["phone"], r.Phone.String)
	require.True(t, r.Phone.Valid)

	require.Equal(t, data["start_date"], r.StartDate.Time.Format("2006-01-02"))
	require.True(t, r.StartDate.Valid)
	require.Equal(t, data["end_date"], r.EndDate.Time.Format("2006-01-02"))
	require.True(t, r.StartDate.Valid)

	require.Equal(t, data["room_id"], fmt.Sprint(r.RoomID))

	require.Equal(t, data["notes"], r.Notes.String)
	require.True(t, r.Notes.Valid)

	require.Equal(t, data["created_at"], r.CreatedAt.Time.Format(time.RFC3339))
	require.True(t, r.CreatedAt.Valid)

	require.Equal(t, data["updated_at"], r.UpdatedAt.Time.Format(time.RFC3339))
	require.True(t, r.UpdatedAt.Valid)

	marshaledData := r.Marshal()
	require.Equal(t, data, marshaledData)
}
