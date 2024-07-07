package yr

import "farstu/internal/cache"

type YRService interface {
	GetForecast() (*YRLocationForecast, error)
}

type yrServiceImpl struct {
	cachedForecast cache.Cache[YRLocationForecast]
}

var _ YRService = yrServiceImpl{}

func (y yrServiceImpl) GetForecast() (*YRLocationForecast, error) {
	return y.cachedForecast.Get()
}

type YRServiceArgs struct {
	ForecastTTL int
	Lat         float64
	Lon         float64
}

func NewYRService(args YRServiceArgs) YRService {
	refreshForecast := func() (*YRLocationForecast, error) {
		return newYRLocationForecast(args.Lat, args.Lon)
	}

	return &yrServiceImpl{
		cachedForecast: cache.New(args.ForecastTTL, refreshForecast),
	}
}
