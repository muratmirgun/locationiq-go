package locationiq

import (
	"context"
)

// SearchLimitCountryURL represents search with limit country url
const SearchLimitCountryURL = "/v1/search.php"

// SearchLimitCountryRequest represents search request
type SearchLimitCountryRequest struct {
	Query        string `query:"q"`
	Format       string `query:"format"`
	CountryCodes string `query:"countrycodes"`
}

// SearchLimitCountryResponse represents search response
type SearchLimitCountryResponse struct {
	PlaceID     string   `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       string   `json:"osm_id"`
	Boundingbox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	Icon        string   `json:"icon,omitempty"`
}

// SearchLimitCountry place
func (c *client) SearchLimitCountry(ctx context.Context, req SearchLimitCountryRequest) ([]SearchLimitCountryResponse, error) {
	url := c.baseURLs[BackendService] + SearchLimitCountryURL + "?key=" + c.Key + "&q=" + req.Query + "&format=" + req.Format + "&countrycodes=" + req.CountryCodes

	var res []SearchLimitCountryResponse

	header := map[string]string{}

	_, err := c.request(ctx, "GET", url, header, &req, &res)

	return res, err
}
