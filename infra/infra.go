package infra

import (
	"time"
)

func StrToTime(t string, tz *time.Location) (time.Time, error) {
	toTime, err := time.ParseInLocation("20060102150405", t, tz)
	if err != nil {
		return time.Time{}, err
	}

	return toTime, nil
}
