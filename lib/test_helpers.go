package lib

import (
	"time"
)

func timeParse(timeString string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeString)
	return t
}
