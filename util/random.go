package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const ABC = "abcdefghijklmnopqrstuvwxyz"

func init() {
	src := rand.NewSource(time.Now().UnixMicro())
	r = rand.New(src)
}

// RandomInt generates a random integer between min and max.
// requirements: min >= 0, max > min. In case of error returns 0.
func RandomInt64(min, max int64) int64 {
	if min < 0 || max < 0 || max <= min {
		return 0
	}

	return min + r.Int63n(max-min+1)
}

// RandomFloat64 generates a random decimal between min and max.
// requirements: min >= 0, max > min. In case of error returns 0.00.
func RandomFloat64(min, max float64) float64 {
	if min < 0.00 || max < 0.00 || max <= min {
		return 0.00
	}

	n := RandomInt64(int64(min*100), int64(max*100))

	return float64(n) / 100
}

// RandomString generates a random string of lenght n.
func RandomString(n int) string {
	var sb strings.Builder
	var b byte

	l := len(ABC)

	for i := 0; i < n; i++ {
		b = ABC[rand.Intn(l)]
		sb.WriteByte(b)
	}

	return sb.String()
}

// RandomDatetime generates a random datetime for the past year.
func RandomDatetime() time.Time {
	standardYear := time.Hour * 24 * 365
	randomDuration := time.Duration(r.Int63n(int64(standardYear)))
	return time.Now().Add(-randomDuration)
}

// RandomName generates a random first or last name of 8 characters long
func RandomName() string {
	return RandomString(8)
}

// RandomName generates a random e-mail of 10 characters followed by @gmail.com
func RandomEmail() string {
	return fmt.Sprint(RandomString(10), "@gmail.com")
}

// RandomPhoneNumber generates a random phone number in the format +000 0000-0000
func RandomPhoneNumber() string {
	return fmt.Sprintf("+%d %d-%d",
		RandomInt64(100, 999),
		RandomInt64(1000, 9999),
		RandomInt64(1000, 9999))
}

// RandomAddress generates a random address
func RandomAddress() string {
	return fmt.Sprintf(
		`street %s %d
		%s %d
		%s`,
		RandomString(8), RandomInt64(10, 99),
		RandomString(8), RandomInt64(10000, 99999),
		RandomString(8))
}

// RandomHourlyFee generates a random hourly fee between 85.00 to 300.00
func RandomHourlyFee() float64 {
	return RandomFloat64(85.00, 300.00)
}

// RandomNote generates a random note
func RandomNote() string {
	return fmt.Sprint("This is a random note:\n",
		RandomString(10), "\n",
		RandomString(10))
}

// RandomLessonDuration returns random int64 between 30 and 240 minutes
func RandomLessonDuration() int64 {
	return RandomInt64(30, 240)
}

// RandomDiscount generates a random discount % between 0.00 (0%) to 0.30 (30%)
func RandomDiscount() float64 {
	return RandomFloat64(0.00, 0.30)
}

// RandomInvoiceAmount returns random float64 between 85.00 and 1200.00
func RandomInvoiceAmount() float64 {
	return RandomFloat64(85.00, 1200.00)
}

// RandomPaymentAmount returns random float64 between 85.00 and 1200.00
func RandomPaymentAmount() float64 {
	return RandomFloat64(85.00, 1200.00)
}
