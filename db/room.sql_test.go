package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createRandomRoom creates a random room in the database
func createRandomRoom(t *testing.T) Room {
	//data := randomRoomData()

	arg := CreateRoomParams{
		Name:          util.RandomName(),
		Description:   util.RandomNote(),
		ImageFilename: fmt.Sprintf("%s.png", util.RandomName()),
	}
	//arg.Unmarshal(data)

	r, err := testStore.CreateRoom(context.Background(), arg)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, arg.Name, r.Name)
	assert.Equal(t, arg.Description, r.Description)
	assert.Equal(t, arg.ImageFilename, r.ImageFilename)
	assert.WithinDuration(t, time.Now(), r.CreatedAt.Time, time.Second)
	assert.True(t, r.CreatedAt.Valid)
	assert.WithinDuration(t, time.Now(), r.UpdatedAt.Time, time.Second)
	assert.True(t, r.UpdatedAt.Valid)

	return r
}

// createRandomRooms create n random rooms in the database
func createRandomRooms(t *testing.T, n int) []Room {
	rooms := make([]Room, n)
	for i := 0; i < n; i++ {
		rooms[i] = createRandomRoom(t)
	}

	return rooms
}

func TestQueries_CreateRoom(t *testing.T) {
	createRandomRoom(t)
}

func TestQueries_ListAvailableRooms(t *testing.T) {
	// remove all restrictions, reservations and rooms
	err := testStore.DeleteAllRoomRestrictions(context.Background())
	require.NoError(t, err)

	err = testStore.DeleteAllReservations(context.Background())
	require.NoError(t, err)

	err = testStore.DeleteAllRooms(context.Background())
	require.NoError(t, err)

	// N is the amount of rooms for the test
	const N = 10

	// startDate is the start date of the series of reservations
	startDate := util.RandomDate()

	rooms := createRandomRooms(t, N)
	reservations := make([]Reservation, N)

	currentStartDate := startDate
	for i, room := range rooms {
		reservations[i] = createRandomWeekReservation(t, room, currentStartDate)
		createRandomRoomRestriction(t, reservations[i])
		currentStartDate = currentStartDate.Add(time.Hour * 24 * 14)
	}

	t.Run("All Rooms Available", func(t *testing.T) {
		arg := ListAvailableRoomsParams{
			Limit:  N * 2,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  startDate.Add(-time.Hour * 24 * 90),
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  startDate.Add(-time.Hour * 24 * 80),
				Valid: true,
			},
		}

		resultRooms, err := testStore.ListAvailableRooms(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, resultRooms, 10)

		for _, v := range resultRooms {
			require.Contains(t, rooms, v)
		}
	})

	t.Run("All Rooms Unavailable", func(t *testing.T) {
		arg := ListAvailableRoomsParams{
			Limit:  N * 2,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  startDate.Add(-time.Hour * 24 * 7),
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  startDate.Add(time.Hour * 24 * 365),
				Valid: true,
			},
		}

		resultRooms, err := testStore.ListAvailableRooms(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, resultRooms, 0)
	})

	t.Run("1st Room Unavailable", func(t *testing.T) {
		arg := ListAvailableRoomsParams{
			Limit:  N * 2,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  startDate.Add(time.Hour * 24 * 2),
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  startDate.Add(time.Hour * 24 * 9),
				Valid: true,
			},
		}

		resultRooms, err := testStore.ListAvailableRooms(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, resultRooms, 9)

		for _, v := range resultRooms {
			require.NotEqual(t, rooms[0].ID, v.ID)
		}
	})

	t.Run("2nd & 3rd Rooms Unavailable", func(t *testing.T) {
		arg := ListAvailableRoomsParams{
			Limit:  N * 2,
			Offset: 0,
			StartDate: pgtype.Date{
				Time:  startDate.Add(time.Hour * 24 * 8),
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  startDate.Add(time.Hour * 24 * 38),
				Valid: true,
			},
		}

		resultRooms, err := testStore.ListAvailableRooms(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, resultRooms, 8)

		for _, v := range resultRooms {
			require.NotEqual(t, rooms[1].ID, v.ID)
			require.NotEqual(t, rooms[2].ID, v.ID)
		}
	})
}
