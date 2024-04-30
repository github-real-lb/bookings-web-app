package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestriction_Scan(t *testing.T) {
	var r Restriction = RestrictionReservation

	err := r.Scan("owner_block")
	require.NoError(t, err)
	require.Equal(t, RestrictionOwnerBlock, r)

	err = r.Scan([]byte("reservation"))
	require.NoError(t, err)
	require.Equal(t, RestrictionReservation, r)

	err = r.Scan(Reservation{})
	require.ErrorContains(t, err, "unsupported scan type for Restriction:")
}
