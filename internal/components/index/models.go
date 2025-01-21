package index

import (
	"farstu/internal/components/clock"
	"farstu/internal/config"
	"farstu/internal/components/sl"
	"farstu/internal/components/yr"
)

type ViewModel struct {
	Config     config.AppConfig
	Departures sl.DeparturesViewModel
	Time       clock.ViewModel
	YRForecast yr.YRForecastViewModel
	YRNow      yr.YRNowViewModel
}
