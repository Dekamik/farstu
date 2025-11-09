package yr

import (
	"log/slog"
	"time"

	"github.com/Dekamik/farstu/internal/cache"
	"github.com/Dekamik/farstu/internal/config"
)

type YRService interface {
	GetViewModels() (YRNowViewModel, YRForecastViewModel)
}

type yrServiceImpl struct {
	appConfig      config.AppConfig
	cachedForecast cache.Cache[yrLocationForecast]
}

var _ YRService = yrServiceImpl{}

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
