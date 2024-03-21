package util

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// N states the number of times to test randomness
const N int = 5

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
			n1 := RandomInt64(test.min, test.max)
			if !test.ok {
				assert.Equal(t, n1, int64(0))
				return
			}

			assert.True(t, n1 >= test.min && n1 <= test.max)

			n2 := RandomInt64(test.min, test.max)
			assert.True(t, n2 >= test.min && n2 <= test.max)

			assert.NotEqual(t, n1, n2)
		})
	}
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n1 := RandomFloat64(test.min, test.max)
			if !test.ok {
				assert.Equal(t, n1, float64(0.00))
				return
			}

			assert.True(t, n1 >= test.min && n1 <= test.max)

			n2 := RandomFloat64(test.min, test.max)
			assert.True(t, n2 >= test.min && n2 <= test.max)

			assert.NotEqual(t, n1, n2)
		})
	}
}

func TestRandomString(t *testing.T) {
	s := RandomString(0)
	assert.Empty(t, s)

	ss := make([]string, N)
	for len := 1; len < 4; len++ {
		for i := 0; i < N; i++ {
			ss[i] = RandomString(len)
			assert.NotEmpty(t, ss[i])
			assert.Len(t, ss[i], len)
		}

		for i := 0; i < N-1; i++ {
			assert.NotEqual(t, ss[i], ss[i+1])
		}
	}
}

func TestRandomDatetime(t *testing.T) {
	dates := make([]time.Time, N)
	for i := 0; i < N; i++ {
		dates[i] = RandomDatetime()
		assert.NotEmpty(t, dates[i])

		days := time.Since(dates[i]).Hours() / 24.00
		assert.True(t, days <= 365.00)
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, dates[i], dates[i+1])
	}
}

func TestRandomName(t *testing.T) {
	names := make([]string, N)
	for i := 0; i < N; i++ {
		names[i] = RandomName()
		assert.NotEmpty(t, names[i])
		assert.Len(t, names[i], 8)
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, names[i], names[i+1])
	}
}

func TestRandomEmail(t *testing.T) {
	emails := make([]string, N)
	for i := 0; i < N; i++ {
		emails[i] = RandomEmail()
		assert.NotEmpty(t, emails[i])
		assert.Len(t, emails[i], 20)
		assert.Contains(t, emails[i], "@gmail.com")
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, emails[i], emails[i+1])
	}
}

func TestRandomPhoneNumber(t *testing.T) {
	phones := make([]string, N)
	for i := 0; i < N; i++ {
		phones[i] = RandomPhoneNumber()
		assert.NotEmpty(t, phones[i])
		assert.Len(t, phones[i], 14)

		assert.Equal(t, phones[i][:1], "+")

		_, err := strconv.Atoi(phones[i][1:4])
		assert.NoError(t, err)

		assert.Equal(t, phones[i][4:5], " ")

		_, err = strconv.Atoi(phones[i][5:9])
		assert.NoError(t, err)

		assert.Equal(t, phones[i][9:10], "-")

		_, err = strconv.Atoi(phones[i][10:14])
		assert.NoError(t, err)
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, phones[i], phones[i+1])
	}
}

func TestRandomAddress(t *testing.T) {
	addresses := make([]string, N)
	for i := 0; i < N; i++ {
		addresses[i] = RandomPhoneNumber()
		assert.NotEmpty(t, addresses[i])
		assert.GreaterOrEqual(t, len(addresses[i]), 10)
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, addresses[i], addresses[i+1])
	}
}

func TestRandomHourlyFee(t *testing.T) {
	fees := make([]float64, N)
	for i := 0; i < N; i++ {
		fees[i] = RandomHourlyFee()
		assert.NotEmpty(t, fees[i])
		assert.GreaterOrEqual(t, fees[i], float64(85.00))
		assert.LessOrEqual(t, fees[i], float64(300.00))
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, fees[i], fees[i+1])
	}
}

func TestRandomNote(t *testing.T) {
	notes := make([]string, N)
	for i := 0; i < N; i++ {
		notes[i] = RandomNote()
		assert.NotEmpty(t, notes[i])
		assert.GreaterOrEqual(t, len(notes[i]), 10)
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, notes[i], notes[i+1])
	}
}

func TestRandomLessonDuration(t *testing.T) {
	durations := make([]int64, N)
	for i := 0; i < N; i++ {
		durations[i] = RandomLessonDuration()
		assert.NotEmpty(t, durations[i])
		assert.GreaterOrEqual(t, durations[i], int64(30))
		assert.LessOrEqual(t, durations[i], int64(240))
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, durations[i], durations[i+1])
	}
}

func TestRandomDiscount(t *testing.T) {
	discounts := make([]float64, N)
	for i := 0; i < N; i++ {
		discounts[i] = RandomDiscount()
		assert.NotEmpty(t, discounts[i])
		assert.GreaterOrEqual(t, discounts[i], float64(0.00))
		assert.LessOrEqual(t, discounts[i], float64(0.30))
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, discounts[i], discounts[i+1])
	}
}

func TestRandomInvoiceAmount(t *testing.T) {
	amounts := make([]float64, N)
	for i := 0; i < N; i++ {
		amounts[i] = RandomHourlyFee()
		assert.NotEmpty(t, amounts[i])
		assert.GreaterOrEqual(t, amounts[i], float64(85.00))
		assert.LessOrEqual(t, amounts[i], float64(1200.00))
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, amounts[i], amounts[i+1])
	}
}

func TestRandomPaymentAmount(t *testing.T) {
	amounts := make([]float64, N)
	for i := 0; i < N; i++ {
		amounts[i] = RandomPaymentAmount()
		assert.NotEmpty(t, amounts[i])
		assert.GreaterOrEqual(t, amounts[i], float64(85.00))
		assert.LessOrEqual(t, amounts[i], float64(1200.00))
	}

	for i := 0; i < N-1; i++ {
		assert.NotEqual(t, amounts[i], amounts[i+1])
	}
}
