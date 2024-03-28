package db

import (
	"context"
	"testing"
	"time"

	"github.com/github-real-lb/bookings-web-app/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomRestriction(t *testing.T) Restriction {
	name := util.RandomName()
	r, err := testStore.CreateRestriction(context.Background(), name)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.Equal(t, name, r.Name)
	assert.WithinDuration(t, time.Now(), r.CreatedAt.Time, time.Second)
	assert.WithinDuration(t, time.Now(), r.UpdatedAt.Time, time.Second)

	return r
}

func TestQueries_CreateRestriction(t *testing.T) {
	createRandomRestriction(t)
}
