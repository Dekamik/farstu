package clock

import (
	"fmt"
	"time"
)

type TimeData struct {
	Time string
	Date string
}

var dayNames = map[time.Weekday]string{
	time.Monday:    "Mån",
	time.Tuesday:   "Tis",
	time.Wednesday: "Ons",
	time.Thursday:  "Tors",
	time.Friday:    "Fre",
	time.Saturday:  "Lör",
	time.Sunday:    "Sön",
}

var monthNames = map[time.Month]string{
	time.January:   "jan",
	time.February:  "feb",
	time.March:     "mar",
	time.April:     "apr",
	time.May:       "maj",
	time.June:      "jun",
	time.July:      "jul",
	time.August:    "aug",
	time.September: "sep",
	time.October:   "okt",
	time.November:  "nov",
	time.December:  "dec",
}

func GetDateStr(date time.Time) string {
	num := date.Format("2")

	var numStr string
	switch rune(num[len(num)-1]) {

	case '1', '2':
		numStr = num + ":a"

	default:
		numStr = num + ":e"
	}

	return fmt.Sprintf("%s %s %s", dayNames[date.Weekday()], numStr, monthNames[date.Month()])
}

func GetTime() TimeData {
	now := time.Now().Local()
	timeStr := now.Format("15:04")
	dateStr := GetDateStr(now)

	return TimeData{
		Time: timeStr,
		Date: dateStr,
	}
}
