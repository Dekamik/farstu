package index

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/yr"
)

type Model struct {
	Config     config.AppConfig
	Time       clock.Model
	YRNow      yr.YRForecastItemModel
	YRForecast yr.YRForecastModel
}
