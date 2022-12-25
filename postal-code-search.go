package locationiq

import (
	"context"
)

// PostalCodeSearchURL represents postal code search url
const PostalCodeSearchURL = "/v1/search.php"

// PostalCodeSearchRequest represents search request
type PostalCodeSearchRequest struct {
	Query          string `query:"q"`
	Format         string `query:"format"`
	PostalCode     string `query:"postalcode"`
	CountryCodes   string `query:"countrycodes"`
	AddressDetails string `query:"addressdetails"`
	Limit          string `query:"limit"`
	NameDetails    string `query:"namedetails"`
	Dedupe         string `query:"dedupe"`
	ViewBox        string `query:"viewbox"`
}

// PostalCodeSearchResponse represents postal code search response
type PostalCodeSearchResponse struct {
	PlaceId     string   `json:"place_id"`
	Licence     string   `json:"licence"`
	Boundingbox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	Address     Address  `json:"address"`
}

type Address struct {
	Suburb      string `json:"suburb"`
	County      string `json:"county"`
	State       string `json:"state"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Town        string `json:"town"`
}

// PostalCodeSearch place
func (c *client) PostalCodeSearch(ctx context.Context, req PostalCodeSearchRequest) ([]PostalCodeSearchResponse, error) {
	url := c.baseURLs[BackendService] + PostalCodeSearchURL + "?key=" +
		c.Key + "&q=" + req.Query +
		"&format=" + req.Format +
		"&countrycodes=" + req.CountryCodes +
		"&postalcode=" + req.PostalCode +
		"&addressdetails=" + req.AddressDetails +
		"&limit=" + req.Limit +
		"&namedetails=" + req.NameDetails +
		"&dedupe=" + req.Dedupe +
		"&viewbox=" + req.ViewBox

	var res []PostalCodeSearchResponse

	header := map[string]string{}

	_, err := c.request(ctx, "GET", url, header, &req, &res)

	return res, err
}
