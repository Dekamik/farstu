package index

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/sl"
	"farstu/internal/yr"
)

type ViewModel struct {
	Config     config.AppConfig
	Departures sl.DeparturesViewModel
	Time       clock.ViewModel
	YRForecast yr.YRForecastViewModel
	YRNow      yr.YRNowViewModel
}
