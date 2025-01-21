package index

import (
	"github.com/Dekamik/farstu/internal/components/clock"
	"github.com/Dekamik/farstu/internal/config"
	"github.com/Dekamik/farstu/internal/components/sl"
	"github.com/Dekamik/farstu/internal/components/yr"
)

type ViewModel struct {
	Config     config.AppConfig
	Departures sl.DeparturesViewModel
	Time       clock.ViewModel
	YRForecast yr.YRForecastViewModel
	YRNow      yr.YRNowViewModel
}
