package yr

import (
	"farstu/internal/config"
	"fmt"
	"time"
)

type YRForecastItemModel struct {
	Enabled            bool
	Time               string
	Precipitation      float64
	PrecipitationColor string
	SymbolCode         string
	SymbolID           string
	Temperature        float64
	TemperatureColor   string
}

func NewYRNowModel(config config.AppConfig, forecast YRLocationForecast) YRForecastItemModel {
	latest := forecast.Properties.Timeseries[0]
	precipitation := latest.Data.Next6Hours.Details.PrecipitationAmount
	symbolCode := latest.Data.Next6Hours.Summary.SymbolCode
	temperature := latest.Data.Instant.Details.AirTemperature

	return YRForecastItemModel{
		Enabled:            config.Weather.Enabled,
		Precipitation:      precipitation,
		PrecipitationColor: getPrecipitationColorClass(config, precipitation),
		SymbolCode:         symbolCode,
		SymbolID:           YRSymbolsID[symbolCode],
		Temperature:        temperature,
		TemperatureColor:   getTemperatureColor(config, temperature),
	}
}

type YRForecastModel struct {
	Enabled bool
	Items   []YRForecastItemModel
}

func NewYRForecastModel(config config.AppConfig, forecast YRLocationForecast) YRForecastModel {
	items := make([]YRForecastItemModel, 0)

	for _, item := range forecast.Properties.Timeseries {
		timeStr := item.Time.Local().Format("15")
		if timeStr != "00" && timeStr != "06" && timeStr != "12" && timeStr != "18" {
			continue
		}

		temperature := item.Data.Instant.Details.AirTemperature
		symbolCode := item.Data.Next6Hours.Summary.SymbolCode
		precipitation := item.Data.Next6Hours.Details.PrecipitationAmount

		forecastItem := YRForecastItemModel{
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

	return YRForecastModel{
		Enabled: config.Weather.Enabled,
		Items:   items[0:config.Weather.MaxRows],
	}
}
