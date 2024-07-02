package views

import "farstu/internal/clock"

type TimeModel struct {
	Date string
	Time string
}

func GetTimeModel() TimeModel {
	t := clock.GetTime()
	return TimeModel{
		Date: t.Date,
		Time: t.Time,
	}
}
