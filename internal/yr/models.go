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

func getForecastAggregates(forecasts yrLocationForecast) []forecastAggregatesResult {
	aggregates := make([]forecastAggregatesResult, 0)
	aggregate := forecastAggregatesResult{}
	isFirst := true
	aggregate.FirstHour = forecasts.Properties.Timeseries[0].Time.Local().Format("15")

	for _, item := range forecasts.Properties.Timeseries {
		hour := item.Time.Local().Format("15")
		if !isFirst && (hour == "00" || hour == "06" || hour == "12" || hour == "18") {
			aggregates = append(aggregates, aggregate)
			aggregate = forecastAggregatesResult{}
			aggregate.FirstHour = hour
		}
		isFirst = false

		aggregate.PrecipitationMin += item.Data.Next1Hours.Details.PrecipitationAmountMin
		aggregate.PrecipitationMax += item.Data.Next1Hours.Details.PrecipitationAmountMax
		aggregate.LastHour = hour

		uvIndex := item.Data.Instant.Details.UltravioletIndexClearSky
		if uvIndex > aggregate.UVIndexMax {
			aggregate.UVIndexMax = uvIndex
		}
	}

	aggregates = append(aggregates, aggregate)
	for i, item := range aggregates {
		aggregates[i].Time = fmt.Sprintf("%s-%s", item.FirstHour, item.LastHour)
	}
	return aggregates
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
	aggregates := getForecastAggregates(forecast)

	for _, item := range forecast.Properties.Timeseries {
		hour := item.Time.Local().Format("15")
		if hour != "00" && hour != "06" && hour != "12" && hour != "18" {
			continue
		}

		temperature := item.Data.Instant.Details.AirTemperature
		symbolCode := item.Data.Next6Hours.Summary.SymbolCode
		precipitationMin := item.Data.Next6Hours.Details.PrecipitationAmountMin
		precipitationMax := item.Data.Next6Hours.Details.PrecipitationAmountMax

		timeStr := fmt.Sprintf("%s-%s", hour, item.Time.Local().Add(time.Hour*6).Format("15"))
		var aggregate forecastAggregatesResult
		for _, item := range aggregates {
			if aggregate.Time == timeStr {
				aggregate = item
				break
			}
		}

		forecastItem := YRForecastItem{
			Time:               timeStr,
			Temperature:        temperature,
			TemperatureColor:   getTemperatureColor(config, temperature),
			SymbolCode:         symbolCode,
			SymbolID:           YRSymbolsID[symbolCode],
			PrecipitationMin:   precipitationMin,
			PrecipitationMax:   precipitationMax,
			PrecipitationColor: getPrecipitationColorClass(config, precipitationMax),
			MaxUVIndex:         aggregate.UVIndexMax,
		}

		forecasts = append(forecasts, forecastItem)
	}

	return YRForecastViewModel{
		Enabled: config.Weather.Enabled,
		Items:   forecasts[0:config.Weather.MaxRows],
	}
}
