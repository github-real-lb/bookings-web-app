package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func TimeToPgDate(time time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  time,
		Valid: true,
	}
}

func TimeToPgTimestamptz(time time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time,
		Valid: true,
	}
}

func Int64ToPgInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: i,
		Valid: true,
	}
}
