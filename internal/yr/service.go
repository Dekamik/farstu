package yr

import (
	"farstu/internal/cache"
	"farstu/internal/config"
	"log/slog"
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
		return newYRLocationForecast(args.Lat, args.Lon)
	}

	return &yrServiceImpl{
		appConfig:      appConfig,
		cachedForecast: cache.New(args.ForecastTTL, refreshForecast),
	}
}
