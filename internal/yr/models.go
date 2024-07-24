package yr

import (
	"farstu/internal/config"
	"fmt"
	"time"
)

type YRForecastItem struct {
	Enabled            bool
	Time               string
	PrecipitationMin   float64
	PrecipitationMax   float64
	PrecipitationColor string
	SymbolCode         string
	SymbolID           string
	Temperature        float64
	TemperatureColor   string
	MaxUVIndex         float64
	UVColor            string
}

type YRNowViewModel struct {
	Enabled  bool
	Forecast *YRForecastItem
	Message  string
}

type forecastAggregatesResult struct {
	Time             string
	FirstHour        string
	LastHour         string
	PrecipitationMin float64
	PrecipitationMax float64
	UVIndexMax       float64
}

func calculateUVColor(uvIndex float64) string {
	if uvIndex < 3 {
		// Green
		return "#3ea72d"
	} else if uvIndex < 6 {
		// Yellow
		return "#fff300"
	} else if uvIndex < 8 {
		// Orange
		return "#f18b00"
	} else if uvIndex < 11 {
		// Red
		return "#e53210"
	} else {
		// Violet
		return "#b567a4"
	}
}

func NewYRNowViewModel(config config.AppConfig, forecast yrLocationForecast) YRNowViewModel {
	latest := forecast.Properties.Timeseries[0]
	precipitationMin := latest.Data.Next6Hours.Details.PrecipitationAmountMin
	precipitationMax := latest.Data.Next6Hours.Details.PrecipitationAmountMax
	symbolCode := latest.Data.Next6Hours.Summary.SymbolCode
	temperature := latest.Data.Instant.Details.AirTemperature
	uvIndex := latest.Data.Instant.Details.UltravioletIndexClearSky

	return YRNowViewModel{
		Enabled: config.Weather.Enabled,
		Forecast: &YRForecastItem{
			Enabled:            config.Weather.Enabled,
			PrecipitationMin:   precipitationMin,
			PrecipitationMax:   precipitationMax,
			PrecipitationColor: getPrecipitationColorClass(config, precipitationMin),
			SymbolCode:         symbolCode,
			SymbolID:           YRSymbolsID[symbolCode],
			Temperature:        temperature,
			TemperatureColor:   getTemperatureColor(config, temperature),
			MaxUVIndex:         uvIndex,
			UVColor:            calculateUVColor(uvIndex),
		},
	}
}

type YRForecastViewModel struct {
	Enabled bool
	Items   []YRForecastItem
	Message string
}

func NewYRForecastViewModel(config config.AppConfig, forecast yrLocationForecast) YRForecastViewModel {
	forecasts := make([]YRForecastItem, 0)

	for i, item := range forecast.Properties.Timeseries {
		hour := item.Time.Local().Format("15")
		if hour != "00" && hour != "06" && hour != "12" && hour != "18" {
			continue
		}

		var largestUVIndex float64 = 0.0
		for j := i; j < i+6; j++ {
			itemUVIndex := forecast.Properties.Timeseries[j].Data.Instant.Details.UltravioletIndexClearSky
			if itemUVIndex > largestUVIndex {
				largestUVIndex = itemUVIndex
			}
		}

		temperature := item.Data.Instant.Details.AirTemperature
		symbolCode := item.Data.Next6Hours.Summary.SymbolCode
		precipitationMin := item.Data.Next6Hours.Details.PrecipitationAmountMin
		precipitationMax := item.Data.Next6Hours.Details.PrecipitationAmountMax

		timeStr := fmt.Sprintf("%s-%s", hour, item.Time.Local().Add(time.Hour*6).Format("15"))
		forecastItem := YRForecastItem{
			Time:               timeStr,
			Temperature:        temperature,
			TemperatureColor:   getTemperatureColor(config, temperature),
			SymbolCode:         symbolCode,
			SymbolID:           YRSymbolsID[symbolCode],
			PrecipitationMin:   precipitationMin,
			PrecipitationMax:   precipitationMax,
			PrecipitationColor: getPrecipitationColorClass(config, precipitationMax),
			MaxUVIndex:         largestUVIndex,
			UVColor:            calculateUVColor(largestUVIndex),
		}

		forecasts = append(forecasts, forecastItem)
	}

	return YRForecastViewModel{
		Enabled: config.Weather.Enabled,
		Items:   forecasts[0:config.Weather.MaxRows],
	}
}
