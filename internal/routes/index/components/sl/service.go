package sl

import (
	"errors"
	"github.com/Dekamik/farstu/internal/cache"
	"github.com/Dekamik/farstu/internal/config"
	"log/slog"
	"strings"
	"time"
)

var ErrSiteIDNotFound = errors.New("site ID not found")

type SLService interface {
	GetDepartures() []Departure
	GetDeviationsViewModel() DeviationsViewModel
}

type slServiceImpl struct {
	appConfig        config.AppConfig
	cachedDepartures cache.Cache[slSiteDeparturesResponse]
	cachedDeviations cache.Cache[[]slDeviationResponse]
}

var _ SLService = slServiceImpl{}

type Departure struct {
	Destination   string
	DisplayTime   string
	LineNumber    string
	TransportMode string
}

func (s slServiceImpl) GetDepartures() []Departure {
	departures := make([]Departure, 0)
	response, err := s.cachedDepartures.Get()

	if err != nil {
		slog.Warn("an error occurred when fetching departures from SL", "err", err)
	} else {
		for _, item := range response.Departures {
			departure := Departure{
				Destination:   item.Destination,
				DisplayTime:   item.Display,
				LineNumber:    item.Line.Designation,
				TransportMode: item.Line.TransportMode,
			}
			departures = append(departures, departure)
		}
	}

	return departures
}

func (s slServiceImpl) GetDeviationsViewModel() DeviationsViewModel {
	var slDeviationsViewModel DeviationsViewModel

	deviations, err := s.cachedDeviations.Get()
	if err != nil {
		slog.Warn("an error occured when fetching deviations from SL", "err", err)
		slDeviationsViewModel = DeviationsViewModel{
			Message: "Fel vid hämtning av störningsinformations",
		}
	} else {
		slDeviationsViewModel = NewDeviationsViewModel(s.appConfig, *deviations)
	}

	return slDeviationsViewModel
}

type SLServiceArgs struct {
	DeparturesTTL  int
	DeviationsTTL  int
	InitRetriesSec []int
	SiteName       string
}

func NewSLService(args SLServiceArgs, appConfig config.AppConfig) (SLService, error) {
	sites, err := callSLSites(false)
	if err != nil {
		if len(args.InitRetriesSec) > 0 {
			for i, retryWaitSecs := range args.InitRetriesSec {
				slog.Info("waiting for next try", "retry", i, "wait_secs", retryWaitSecs)
				time.Sleep(time.Duration(retryWaitSecs) * time.Second)

				sites, err = callSLSites(false)
				if err == nil {
					break
				}
			}

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var siteID int
	queryString := strings.ToLower(args.SiteName)
	siteIDFound := false
	for _, site := range *sites {
		if strings.ToLower(site.Name) == queryString {
			siteID = site.ID
			siteIDFound = true
			break
		}
	}

	if !siteIDFound {
		return nil, ErrSiteIDNotFound
	}

	refreshDepartures := func() (*slSiteDeparturesResponse, error) {
		return callSLSiteDepartures(siteID)
	}

	refreshDeviations := func() (*[]slDeviationResponse, error) {
		deviationsArgs := callSLDeviationsArgs{
			Future: appConfig.SL.Deviations.Future,
			Lines:  appConfig.SL.Deviations.Lines,
			Sites:  appConfig.SL.Deviations.Sites,
		}
		return callSLDeviations(deviationsArgs)
	}

	return slServiceImpl{
		appConfig:        appConfig,
		cachedDepartures: cache.New(args.DeparturesTTL, refreshDepartures),
		cachedDeviations: cache.New(args.DeviationsTTL, refreshDeviations),
	}, nil
}
