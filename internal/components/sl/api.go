package sl

import (
	"fmt"
	"time"

	"github.com/Dekamik/farstu/internal/api"
)

type slSitesResponseItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type slSiteDeparturesResponse struct {
	Departures []struct {
		Destination string `json:"destination"`
		Display     string `json:"display"`
		Line        struct {
			Designation   string `json:"designation"`
			TransportMode string `json:"transport_mode"`
		} `json:"line"`
	} `json:"departures"`
}

type slDeviationsResponse struct {
	Deviations []struct {
		Publish struct {
			From time.Time `json:"from"`
			Upto time.Time `json:"upto"`
		} `json:"publish"`

		Priority []struct {
			ImportanceLevel int `json:"importance_level"`
			InfluenceLevel int `json:"influence_level"`
			UrgencyLevel int `json:"urgency_level"`
		} `json:"priority"`

		MessageVariants []struct {
			Header string `json:"header"`
			Details string `json:"details"`
			ScopeAlias string `json:"scope_alias"`
			Weblink string `json:"weblink"`
			Language string `json:"language"`
		}
		Lines []struct {
			ID int `json:"id"`
			Designation string `json:"designation"`
			TransportMode string `json:"transportMode"`
			Name string `json:"name"`
			GroupOfLines string `json:"group_of_lines"`
		}
	}
}

func callSLSites(expand bool) (*[]slSitesResponseItem, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites?expand=%t", expand)
	return api.GET[[]slSitesResponseItem](url)
}

func callSLSiteDepartures(siteID int) (*slSiteDeparturesResponse, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites/%d/departures", siteID)
	return api.GET[slSiteDeparturesResponse](url)
}

type callSLDeviationsArgs struct {
	Future bool
	TransportAuthority *int
	Sites []int
	Lines []int
	TransportMode string
}

func callSLDeviations(args callSLDeviationsArgs) (*slDeviationsResponse, error) {
	if args.TransportAuthority == nil {
		default_ta := 1
		args.TransportAuthority = &default_ta
	}

	url := fmt.Sprintf("https://deviations.integration.sl.se/v1/messages?future=%t&transport_authority=%d", args.Future, args.TransportAuthority)

	for _, site := range args.Sites {
		url = fmt.Sprintf("%s&site=%d", url, site)
	}
	for _, line := range args.Lines {
		url = fmt.Sprintf("%s&line=%d", url, line)
	}
	if args.TransportMode != "" {
		url = fmt.Sprintf("%s&transport_mode=%s", url, args.TransportMode)
	}

	return api.GET[slDeviationsResponse](url)
}
