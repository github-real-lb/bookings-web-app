package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type KeysMap map[any]bool

// RequireUnique checks if value exists in values.
// If value exists it fails the test, otherwise it adds value to values.
func RequireUnique(t *testing.T, value any, values KeysMap) {
	require.NotNil(t, value)
	require.NotNil(t, values)

	if len(values) != 0 {
		_, exist := values[value]
		if exist {
			require.Falsef(t, exist, `value "%v" is not unique`, value)
			// t.Errorf(`value "%v" is not unique`, value)
			// return
		}
	}

	values[value] = true
}
