package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomRoomData() map[string]string {
	data := make(map[string]string)
	data["name"] = util.RandomName()
	data["description"] = util.RandomNote()
	data["image_filename"] = fmt.Sprintf("%s.png", util.RandomName())
	return data
}

func createRandomRoom(t *testing.T) Room {
	data := randomRoomData()

	arg := CreateRoomParams{}
	arg.Unmarshal(data)

	r, err := testStore.CreateRoom(context.Background(), arg)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, arg.Name, r.Name)
	assert.Equal(t, arg.Description, r.Description)
	assert.Equal(t, arg.ImageFilename, r.ImageFilename)
	assert.WithinDuration(t, time.Now(), r.CreatedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now(), r.UpdatedAt.Time, time.Second)

	return r
}

func TestQueries_CreateRoom(t *testing.T) {
	createRandomRoom(t)
}
