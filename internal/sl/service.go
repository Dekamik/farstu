package sl

import (
	"errors"
	"farstu/internal/config"
	"log/slog"
	"strings"
	"time"
)

type SLService interface {
	GetViewModel() DeparturesViewModel
}

type slServiceImpl struct {
	appConfig config.AppConfig
	siteID    int
}

var _ SLService = slServiceImpl{}

func (s slServiceImpl) GetViewModel() DeparturesViewModel {
	var slDeparturesViewModel DeparturesViewModel

	departures, err := getSLSiteDepartures(s.siteID)
	if err != nil {
		slog.Warn("an error occurred when fetching departures from SL", "err", err)
		slDeparturesViewModel = DeparturesViewModel{
			Enabled: s.appConfig.SL.Enabled,
			Message: "Fel vid avgångsdatahämtning",
		}
	} else {
		slDeparturesViewModel = NewDeparturesViewModel(s.appConfig, *departures)
	}

	return slDeparturesViewModel
}

var ErrSiteIDNotFound = errors.New("site ID not found")

type SLServiceArgs struct {
	DeparturesTTL  int
	InitRetriesSec []int
	SiteName       string
}

func NewSLService(args SLServiceArgs, appConfig config.AppConfig) (SLService, error) {
	sites, err := getSLSites(false)
	if err != nil {
		if len(args.InitRetriesSec) > 0 {
			for i, retryWaitSecs := range args.InitRetriesSec {
				slog.Info("waiting for next try", "retry", i, "wait_secs", retryWaitSecs)
				time.Sleep(time.Duration(retryWaitSecs) * time.Second)

				sites, err = getSLSites(false)
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

	return slServiceImpl{
		appConfig: appConfig,
		siteID: siteID,
	}, nil
}
