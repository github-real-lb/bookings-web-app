package util

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// N states the number of times to test randomness
const N int = 5

type keysMap map[any]bool

// requireUnique checks if key exists in keys.
// If key exists it fails the test, otherwise it adds the key to the map.
func requireUnique(t *testing.T, key any, keys keysMap) {
	require.NotNil(t, key)
	require.NotNil(t, keys)

	if len(keys) != 0 {
		_, exist := keys[key]
		require.False(t, exist)
	}

	keys[key] = true
}

func TestRandomInt64(t *testing.T) {
	tests := []struct {
		name string
		min  int64
		max  int64
		ok   bool
	}{
		{name: "OK", min: 0, max: 100, ok: true},
		{name: "Negative min", min: -10, max: 10, ok: false},
		{name: "Negative max", min: -100, max: -10, ok: false},
		{name: "max < min", min: 100, max: 10, ok: false},
		{name: "max = min", min: 100, max: 100, ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := RandomInt64(test.min, test.max)
			if test.ok {
				assert.True(t, n >= test.min && n <= test.max)
			} else {
				assert.Equal(t, n, int64(0))
			}
		})
	}

	t.Run("Randomness", func(t *testing.T) {
		numbers := make(keysMap)
		min := int64(0)
		max := int64(50)
		for i := 0; i < N; i++ {
			number := RandomInt64(min, max)
			assert.True(t, number >= min && number <= max)

			requireUnique(t, number, numbers)
		}
	})
}

func TestRandomFloat64(t *testing.T) {
	tests := []struct {
		name string
		min  float64
		max  float64
		ok   bool
	}{
		{name: "OK 0.00 -> 100.00", min: 0.00, max: 100.00, ok: true},
		{name: "OK 0.00 -> 1.00", min: 0.00, max: 1.00, ok: true},
		{name: "OK 0.40 -> 0.60", min: 0.40, max: 0.60, ok: true},
		{name: "Negative min", min: -10.00, max: 10.00, ok: false},
		{name: "Negative max", min: -100.00, max: -10.00, ok: false},
		{name: "max < min", min: 100.00, max: 10.00, ok: false},
		{name: "max = min", min: 10.00, max: 10.00, ok: false},
	}

	// test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := RandomFloat64(test.min, test.max)

			if test.ok {
				assert.True(t, n >= test.min && n <= test.max)
			} else {
				assert.Equal(t, n, float64(0.00))
			}
		})
	}

	t.Run("Randomness", func(t *testing.T) {
		numbers := make(keysMap)
		min := 0.00
		max := 50.00
		for i := 0; i < N; i++ {
			number := RandomFloat64(min, max)
			assert.True(t, number >= min && number <= max)

			requireUnique(t, number, numbers)
		}
	})
}

func TestRandomString(t *testing.T) {
	s := RandomString(0)
	assert.Empty(t, s)

	ss := make(keysMap)
	for len := 2; len < 4; len++ {
		for i := 0; i < N; i++ {
			s := RandomString(len)
			assert.NotEmpty(t, s)
			assert.Len(t, s, len)

			requireUnique(t, s, ss)
		}
	}
}

func TestRandomDate(t *testing.T) {
	dates := make(keysMap)
	for i := 0; i < N; i++ {
		date := RandomDate()

		days := time.Since(date).Hours() / 24.00
		assert.True(t, days <= 365.00)
		assert.Equal(t, 0, date.Hour())
		assert.Equal(t, 0, date.Minute())
		assert.Equal(t, 0, date.Second())

		requireUnique(t, date, dates)

	}
}

func TestRandomDatetime(t *testing.T) {
	dates := make(keysMap)
	for i := 0; i < N; i++ {
		date := RandomDatetime()

		days := time.Since(date).Hours() / 24.00
		assert.True(t, days <= 365.00)

		requireUnique(t, date, dates)
	}
}

func TestRandomName(t *testing.T) {
	names := make(keysMap)
	for i := 0; i < N; i++ {
		name := RandomName()
		assert.NotEmpty(t, name)
		assert.Len(t, name, 8)

		requireUnique(t, name, names)
	}

}

func TestRandomEmail(t *testing.T) {
	emails := make(keysMap)
	for i := 0; i < N; i++ {
		email := RandomEmail()
		assert.NotEmpty(t, email)
		assert.Len(t, email, 20)
		assert.Contains(t, email, "@gmail.com")

		requireUnique(t, email, emails)
	}
}

func TestRandomPhoneNumber(t *testing.T) {
	phones := make(keysMap)
	for i := 0; i < N; i++ {
		phone := RandomPhoneNumber()
		assert.NotEmpty(t, phone)
		assert.Len(t, phone, 14)

		assert.Equal(t, phone[:1], "+")

		_, err := strconv.Atoi(phone[1:4])
		assert.NoError(t, err)

		assert.Equal(t, phone[4:5], " ")

		_, err = strconv.Atoi(phone[5:9])
		assert.NoError(t, err)

		assert.Equal(t, phone[9:10], "-")

		_, err = strconv.Atoi(phone[10:14])
		assert.NoError(t, err)

		requireUnique(t, phone, phones)
	}
}

func TestRandomAddress(t *testing.T) {
	addresses := make(keysMap)
	for i := 0; i < N; i++ {
		address := RandomAddress()
		assert.NotEmpty(t, address)
		assert.GreaterOrEqual(t, len(address), 10)

		requireUnique(t, address, addresses)
	}

}

func TestRandomHourlyFee(t *testing.T) {
	fees := make(keysMap)
	for i := 0; i < N; i++ {
		fee := RandomHourlyFee()
		assert.GreaterOrEqual(t, fee, float64(85.00))
		assert.LessOrEqual(t, fee, float64(300.00))

		requireUnique(t, fee, fees)
	}

}

func TestRandomNote(t *testing.T) {
	notes := make(keysMap)
	for i := 0; i < N; i++ {
		note := RandomNote()
		assert.NotEmpty(t, note)
		assert.GreaterOrEqual(t, len(note), 10)

		requireUnique(t, note, notes)
	}

}

func TestRandomLessonDuration(t *testing.T) {
	durations := make(keysMap)
	for i := 0; i < N; i++ {
		duration := RandomLessonDuration()
		assert.GreaterOrEqual(t, duration, int64(30))
		assert.LessOrEqual(t, duration, int64(240))

		requireUnique(t, duration, durations)
	}
}

func TestRandomDiscount(t *testing.T) {
	discounts := make(keysMap)
	for i := 0; i < N; i++ {
		discount := RandomDiscount()
		assert.GreaterOrEqual(t, discount, float64(0.00))
		assert.LessOrEqual(t, discount, float64(0.30))

		requireUnique(t, discount, discounts)
	}
}

func TestRandomInvoiceAmount(t *testing.T) {
	amounts := make(keysMap)
	for i := 0; i < N; i++ {
		amount := RandomInvoiceAmount()
		assert.GreaterOrEqual(t, amount, float64(85.00))
		assert.LessOrEqual(t, amount, float64(1200.00))

		requireUnique(t, amount, amounts)
	}
}

func TestRandomPaymentAmount(t *testing.T) {
	amounts := make(keysMap)
	for i := 0; i < N; i++ {
		amount := RandomPaymentAmount()
		assert.GreaterOrEqual(t, amount, float64(85.00))
		assert.LessOrEqual(t, amount, float64(1200.00))

		requireUnique(t, amount, amounts)
	}
}

func TestGenerateReservationCode(t *testing.T) {
	reservationName := RandomString(4)
	code, err := GenerateReservationCode(reservationName, 5)
	assert.Empty(t, code)
	require.Error(t, err)

	codes := make(keysMap)
	for i := 0; i < N; i++ {
		code, err := GenerateReservationCode(reservationName, 6)
		assert.Len(t, code, 6)
		require.NoError(t, err)

		requireUnique(t, code, codes)
	}
}
