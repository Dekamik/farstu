package index

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/gtfs"
	"farstu/internal/yr"
)

type ViewModel struct {
	Config     config.AppConfig
	Departures gtfs.DeparturesViewModel
	Time       clock.ViewModel
	YRNow      yr.YRNowViewModel
	YRForecast yr.YRForecastViewModel
}
