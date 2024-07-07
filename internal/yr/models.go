package yr

import (
	"farstu/internal/config"
	"fmt"
	"time"
)

type YRForecastItem struct {
	Enabled            bool
	Time               string
	Precipitation      float64
	PrecipitationColor string
	SymbolCode         string
	SymbolID           string
	Temperature        float64
	TemperatureColor   string
}

type YRNowViewModel struct {
	Enabled  bool
	Forecast *YRForecastItem
	Message  string
}

func NewYRNowViewModel(config config.AppConfig, forecast YRLocationForecast) YRNowViewModel {
	latest := forecast.Properties.Timeseries[0]
	precipitation := latest.Data.Next6Hours.Details.PrecipitationAmount
	symbolCode := latest.Data.Next6Hours.Summary.SymbolCode
	temperature := latest.Data.Instant.Details.AirTemperature

	return YRNowViewModel{
		Enabled: config.Weather.Enabled,
		Forecast: &YRForecastItem{
			Enabled:            config.Weather.Enabled,
			Precipitation:      precipitation,
			PrecipitationColor: getPrecipitationColorClass(config, precipitation),
			SymbolCode:         symbolCode,
			SymbolID:           YRSymbolsID[symbolCode],
			Temperature:        temperature,
			TemperatureColor:   getTemperatureColor(config, temperature),
		},
	}
}

type YRForecastViewModel struct {
	Enabled bool
	Items   []YRForecastItem
	Message string
}

func NewYRForecastViewModel(config config.AppConfig, forecast YRLocationForecast) YRForecastViewModel {
	items := make([]YRForecastItem, 0)

	for _, item := range forecast.Properties.Timeseries {
		timeStr := item.Time.Local().Format("15")
		if timeStr != "00" && timeStr != "06" && timeStr != "12" && timeStr != "18" {
			continue
		}

		temperature := item.Data.Instant.Details.AirTemperature
		symbolCode := item.Data.Next6Hours.Summary.SymbolCode
		precipitation := item.Data.Next6Hours.Details.PrecipitationAmount

		forecastItem := YRForecastItem{
			Time:               fmt.Sprintf("%s-%s", timeStr, item.Time.Local().Add(time.Hour*6).Format("15")),
			Temperature:        temperature,
			TemperatureColor:   getTemperatureColor(config, temperature),
			SymbolCode:         symbolCode,
			SymbolID:           YRSymbolsID[symbolCode],
			Precipitation:      precipitation,
			PrecipitationColor: getPrecipitationColorClass(config, precipitation),
		}

		items = append(items, forecastItem)
	}

	return YRForecastViewModel{
		Enabled: config.Weather.Enabled,
		Items:   items[0:config.Weather.MaxRows],
	}
}
