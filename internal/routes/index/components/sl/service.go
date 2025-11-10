package sl

import (
	"errors"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/Dekamik/farstu/internal/cache"
	"github.com/Dekamik/farstu/internal/config"
)

var ErrSiteIDNotFound = errors.New("site ID not found")

type SLService interface {
	GetDepartures() []Departure
	GetDeviations() []Deviation
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

	if err == nil {
		for _, item := range response.Departures {
			departure := Departure{
				Destination:   item.Destination,
				DisplayTime:   item.Display,
				LineNumber:    item.Line.Designation,
				TransportMode: item.Line.TransportMode,
			}
			departures = append(departures, departure)
		}
	} else {
		slog.Warn("an error occurred when fetching departures from SL", "err", err)
	}

	return departures
}

func (s slServiceImpl) GetDeviations() []Deviation {
	deviations := make([]Deviation, 0)
	response, err := s.cachedDeviations.Get()

	if err == nil {
		for i, item := range (*response) {
			deviation := Deviation{
				Priority: DeviationPriority{
					ImportanceLevel: item.Priority.ImportanceLevel,
					InfluenceLevel:  item.Priority.InfluenceLevel,
					UrgencyLevel:    item.Priority.UrgencyLevel,
				},
				MessageVariants: make(map[string]DeviationMessage),
				Lines:           make([]DeviationLine, 0),
			}

			for _, message := range (*response)[i].MessageVariants {
				m := DeviationMessage{
					Header:     message.Header,
					Details:    message.Details,
					ScopeAlias: message.ScopeAlias,
					Weblink:    message.Weblink,
				}
				deviation.MessageVariants[message.Language] = m
			}

			for _, line := range (*response)[i].Scope.Lines {
				l := DeviationLine{
					ID:            line.ID,
					Designation:   line.Designation,
					TransportMode: line.TransportMode,
					Name:          line.Name,
					GroupOfLines:  line.GroupOfLines,
				}
				deviation.Lines = append(deviation.Lines, l)
			}

			deviation.Render = calculateRender(deviation)

			deviations = append(deviations, deviation)
		}

		sort.Slice(deviations, func(i, j int) bool {
			return deviations[i].Priority.UrgencyLevel > deviations[j].Priority.UrgencyLevel
		})
		sort.Slice(deviations, func(i, j int) bool {
			return deviations[i].Priority.InfluenceLevel > deviations[j].Priority.InfluenceLevel
		})
		sort.Slice(deviations, func(i, j int) bool {
			return deviations[i].Priority.ImportanceLevel > deviations[j].Priority.ImportanceLevel
		})
	} else {
		slog.Warn("an error occurred when fetching deviations from SL", "err", err)
	}

	return deviations
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
