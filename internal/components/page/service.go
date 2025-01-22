package page

import (
	"github.com/Dekamik/farstu/internal/components/yr"
	"github.com/Dekamik/farstu/internal/config"
	"time"
)

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
	response, err := yr.CallYRSunrise(appConfig.Weather.Lat, appConfig.Weather.Lon, time.Now())
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
		response, err := yr.CallYRSunrise(appConfig.Weather.Lat, appConfig.Weather.Lon, tomorrow)
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
