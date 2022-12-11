package locationiq

import (
	"context"
)

// SearchURL represents search url
const SearchURL = "/v1/search.php"

// SearchRequest represents search request
type SearchRequest struct {
	Key    string `query:"key"`
	Query  string `query:"q"`
	Format string `query:"format"`
}

// SearchResponse represents search response
type SearchResponse struct {
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
	Icon        string   `json:"icon"`
}

// Search place
func (c *client) Search(ctx context.Context, req SearchRequest) ([]SearchResponse, error) {
	url := c.baseURLs[BackendService] + SearchURL + "?key=" + req.Key + "&q=" + req.Query + "&format=" + req.Format

	var res []SearchResponse

	header := map[string]string{}

	_, err := c.request(ctx, "GET", url, header, &req, &res)

	return res, err
}
