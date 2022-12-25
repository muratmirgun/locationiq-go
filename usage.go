package locationiq

import (
	"context"
)

// BalanceURL represents balance general usage url
const BalanceURL = "/v1/search.php"

// BalanceRequest represents balance and general usage request
type BalanceRequest struct{}

// BalanceResponse represents balance and general usage response
type BalanceResponse struct {
	Status  string  `json:"status"`
	Balance Balance `json:"balance"`
}

type Balance struct {
	Day   int `json:"day"`
	Bonus int `json:"bonus"`
}

// Balance general usage
func (c *client) Balance(ctx context.Context, req BalanceRequest) (BalanceResponse, error) {
	url := c.baseURLs[BackendService] + BalanceURL + "?key=" + c.Key

	var res BalanceResponse

	header := map[string]string{}

	_, err := c.request(ctx, "GET", url, header, &req, &res)

	return res, err
}
