package utils

import (
	"time"
)

func ConvertToIstanbulTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		return t
	}
	return t.In(loc)
}
