package main

import (
	"strings"
	"time"
)

// SBool is a bool that can be parsed from JSON as Ja or Nej
type SBool bool

func (b *SBool) UnmarshalJSON(data []byte) error {
	str := string(data)
	*b = strings.ToLower(str) == "ja"
	return nil
}

type SHolidayResponse struct {
	CacheTime time.Time `json:"cachetid"`
	Version   string    `json:"version"`
	URI       string    `json:"uri"`
	StartDate time.Time `json:"startdatum"`
	EndDate   time.Time `json:"slutdatum"`
	Days      struct {
		Date         time.Time `json:"datum"`
		Weekday      string    `json:"veckodag"`
		IsNotWorkday SBool     `json:"arbetsfri dag"`
		IsRedDay     SBool     `json:"röd dag"`
		WeekNumber   string    `json:"vecka"`
		DayInWeek    string    `json:"dag i vecka"`
		Holiday      string    `json:"helgdag"`
		NamesDay     []string  `json:"namnsdag"`
		FlagDay      string    `json:"flaggdag"`
	} `json:"dagar"`
}

func main() {
	// TODO: Scrape this API: https://api.dryg.net/ to gather Swedish holidays
	// Exempel: https://sholiday.faboul.se/dagar/v2.1/2024
}
