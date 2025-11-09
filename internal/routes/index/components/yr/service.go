package yr

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Dekamik/farstu/internal/cache"
	"github.com/Dekamik/farstu/internal/config"
	"github.com/Dekamik/farstu/internal/routes/shared"
)

type YRService interface {
	GetForecast(config.AppConfig) []YRForecastItem
	GetViewModels() (YRNowViewModel, YRForecastViewModel)
}

type yrServiceImpl struct {
	appConfig      config.AppConfig
	cachedForecast cache.Cache[yrLocationForecast]
}

var _ YRService = yrServiceImpl{}

func (y yrServiceImpl) GetForecast(config config.AppConfig) []YRForecastItem {
	forecasts := make([]YRForecastItem, 0)
	response, err := y.cachedForecast.Get()

	if err != nil {
		slog.Warn("an error occurred when fetching weather forecasts from YR.no", "err", err)
	} else {
		for i, item := range response.Properties.Timeseries {
			hour := item.Time.Local().Format("15")
			if i != 0 && hour != "00" && hour != "06" && hour != "12" && hour != "18" {
				continue
			}

			var largestUVIndex float64 = 0.0
			for j := i; j < i+6; j++ {
				itemUVIndex := response.Properties.Timeseries[j].Data.Instant.Details.UltravioletIndexClearSky
				if itemUVIndex > largestUVIndex {
					largestUVIndex = itemUVIndex
				}
			}

			temperature := item.Data.Instant.Details.AirTemperature
			symbolCode := item.Data.Next6Hours.Summary.SymbolCode
			precipitationMin := item.Data.Next6Hours.Details.PrecipitationAmountMin
			precipitationMax := item.Data.Next6Hours.Details.PrecipitationAmountMax

			weekDay := item.Time.Local().Weekday()
			timeStr := fmt.Sprintf("%s %s-%s", shared.DayNames[weekDay], hour, item.Time.Local().Add(time.Hour*6).Format("15"))

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
	}

	return forecasts
}

func (y yrServiceImpl) GetViewModels() (YRNowViewModel, YRForecastViewModel) {
	var yrNowViewModel YRNowViewModel
	var yrForecastViewModel YRForecastViewModel

	forecast, err := y.cachedForecast.Get()
	if err != nil {
		slog.Warn("an error occurred when fetching weather forcasts from YR.no", "err", err)
		yrNowViewModel = YRNowViewModel{
			Enabled: y.appConfig.Weather.Enabled,
			Message: "Fel vid väderdatahämtning",
		}
		yrForecastViewModel = YRForecastViewModel{
			Enabled: y.appConfig.Weather.Enabled,
			Message: "Fel vid väderdatahämtning",
		}
	} else {
		yrNowViewModel = NewYRNowViewModel(y.appConfig, *forecast)
		yrForecastViewModel = NewYRForecastViewModel(y.appConfig, *forecast)
	}

	return yrNowViewModel, yrForecastViewModel
}

type YRServiceArgs struct {
	ForecastTTL int
	Lat         float64
	Lon         float64
}

func NewYRService(args YRServiceArgs, appConfig config.AppConfig) YRService {
	refreshForecast := func() (*yrLocationForecast, error) {
		return callYRLocationForecast(args.Lat, args.Lon)
	}

	return &yrServiceImpl{
		appConfig:      appConfig,
		cachedForecast: cache.New(args.ForecastTTL, refreshForecast),
	}
}

var monthToSeason = map[string]string{
	"01": "winter",
	"02": "winter",
	"03": "spring",
	"04": "spring",
	"05": "spring",
	"06": "summer",
	"07": "summer",
	"08": "summer",
	"09": "autumn",
	"10": "autumn",
	"11": "autumn",
	"12": "winter",
}

type SeasonAndTimeOfDay struct {
	Season                   string
	SecondsUntilNextSunEvent int
	TimeOfDay                string
}

func GetSeasonAndTimeOfDay(appConfig config.AppConfig) (*SeasonAndTimeOfDay, error) {
	response, err := CallYRSunrise(appConfig.Weather.Lat, appConfig.Weather.Lon, time.Now())
	if err != nil {
		return nil, err
	}

	month := time.Now().Format("01")
	season := monthToSeason[month]

	var timeOfDay string
	now := time.Now().Local()
	sunrise := response.Properties.Sunrise.Time.Local()
	sunset := response.Properties.Sunset.Time.Local()
	fiveMinutesBeforeSunrise := sunrise.Add(time.Duration(-5) * time.Minute)
	fiveMinutesBeforeSunset := sunset.Add(time.Duration(-5) * time.Minute)
	if now.After(fiveMinutesBeforeSunrise) && now.Before(fiveMinutesBeforeSunset) {
		timeOfDay = "day"
	} else {
		timeOfDay = "night"
	}

	var secondsUntilNextSunriseOrSunset float64
	if now.Before(sunrise) {
		secondsUntilNextSunriseOrSunset = sunrise.Sub(now).Seconds()
	} else if now.Before(sunset) {
		secondsUntilNextSunriseOrSunset = sunset.Sub(now).Seconds()
	} else {
		// Next event is tomorrows sunrise
		tomorrow := now.Add(time.Duration(24) * time.Hour)
		response, err := CallYRSunrise(appConfig.Weather.Lat, appConfig.Weather.Lon, tomorrow)
		if err != nil {
			return nil, err
		}
		secondsUntilNextSunriseOrSunset = response.Properties.Sunrise.Time.Local().Sub(now).Seconds()
	}

	return &SeasonAndTimeOfDay{
		Season:                   season,
		SecondsUntilNextSunEvent: int(secondsUntilNextSunriseOrSunset),
		TimeOfDay:                timeOfDay,
	}, nil
}
