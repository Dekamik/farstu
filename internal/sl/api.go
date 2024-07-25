package sl

import (
	"farstu/internal/api"
	"fmt"
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

func callSLSites(expand bool) (*[]slSitesResponseItem, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites?expand=%t", expand)
	return api.GET[[]slSitesResponseItem](url)
}

func callSLSiteDepartures(siteID int) (*slSiteDeparturesResponse, error) {
	url := fmt.Sprintf("https://transport.integration.sl.se/v1/sites/%d/departures", siteID)
	return api.GET[slSiteDeparturesResponse](url)
}
