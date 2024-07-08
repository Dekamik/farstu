package yr

import "farstu/internal/cache"

type YRService interface {
	GetForecast() (*yrLocationForecast, error)
}

type yrServiceImpl struct {
	cachedForecast cache.Cache[yrLocationForecast]
}

var _ YRService = yrServiceImpl{}

func (y yrServiceImpl) GetForecast() (*yrLocationForecast, error) {
	return y.cachedForecast.Get()
}

type YRServiceArgs struct {
	ForecastTTL int
	Lat         float64
	Lon         float64
}

func NewYRService(args YRServiceArgs) YRService {
	refreshForecast := func() (*yrLocationForecast, error) {
		return newYRLocationForecast(args.Lat, args.Lon)
	}

	return &yrServiceImpl{
		cachedForecast: cache.New(args.ForecastTTL, refreshForecast),
	}
}
