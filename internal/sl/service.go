package sl

import (
	"errors"
	"farstu/internal/cache"
	"log/slog"
	"strings"
	"time"
)

type SLService interface {
	GetDepartures() (*slSiteDeparturesResponse, error)
}

type slServiceImpl struct {
	cachedDepartures cache.Cache[slSiteDeparturesResponse]
}

var _ SLService = slServiceImpl{}

func (s slServiceImpl) GetDepartures() (*slSiteDeparturesResponse, error) {
	return s.cachedDepartures.Get()
}

var ErrSiteIDNotFound = errors.New("site ID not found")

type SLServiceArgs struct {
	DeparturesTTL int
	RetriesSec       []int
	SiteName      string
}

func NewSLService(args SLServiceArgs) (SLService, error) {
	sites, err := getSLSites(false)
	if err != nil {
		if len(args.RetriesSec) > 0 {
			for i, retryWaitSecs := range args.RetriesSec {
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

	refreshDepartures := func() (*slSiteDeparturesResponse, error) {
		return getSLSiteDepartures(siteID)
	}

	return slServiceImpl{
		cachedDepartures: cache.New(args.DeparturesTTL, refreshDepartures),
	}, nil
}
