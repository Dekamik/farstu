package sl

import (
	"farstu/internal/api"
	"fmt"
)

type SLSitesResponseItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SLSiteDeparturesResponse struct {
	Departures []struct {
		Destination string `json:"destination"`
		Display     string `json:"display"`
		Line        struct {
			Designation   string `json:"designation"`
			TransportMode string `json:"transport_mode"`
		} `json:"line"`
	} `json:"departures"`
}

func getSLSites(expand bool) (*[]SLSitesResponseItem, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites?expand=%t", expand)
	return api.GET[[]SLSitesResponseItem](url)
}

func getSLSiteDepartures(siteID int) (*SLSiteDeparturesResponse, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites/%d/departures", siteID)
	return api.GET[SLSiteDeparturesResponse](url)
}
