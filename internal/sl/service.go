package sl

import (
	"errors"
	"farstu/internal/cache"
	"strings"
)

type SLService interface {
	GetDepartures() (*SLSiteDeparturesResponse, error)
}

type slServiceImpl struct{
	cachedDepartures cache.Cache[SLSiteDeparturesResponse]
}

var _ SLService = slServiceImpl{}

func (s slServiceImpl) GetDepartures() (*SLSiteDeparturesResponse, error) {
	return s.cachedDepartures.Get()
}

var ErrSiteIDNotFound = errors.New("site ID not found")

type SLServiceArgs struct {
	DeparturesTTL int
	SiteName      string
}

func NewSLService(args SLServiceArgs) (SLService, error) {
	sites, err := getSLSites(false)
	if err != nil {
		return nil, err
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

	refreshDepartures := func() (*SLSiteDeparturesResponse, error) {
		return getSLSiteDepartures(siteID)
	}
	
	return slServiceImpl{
		cachedDepartures: cache.New(args.DeparturesTTL, refreshDepartures),
	}, nil
}