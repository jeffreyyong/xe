package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseEndpoint   = "https://api.exchangeratesapi.io"
	LatestPath     = "latest"
	HistoricalPath = "history"
	SymbolsParam   = "symbols"
	BaseParam      = "base"
	StartDateParam = "start_at"
	EndDateParam   = "end_at"
	EURSymbol      = "EUR"
)

var (
	// RetryCount specifies the number of retries a resty
	// client will perform after a failed request.
	RetryCount = 3

	// RetryWaitTime specifies the wait time the client
	// will wait before retrying
	RetryWaitTime = 500 * time.Millisecond

	// RetryMaxWaitTime specifies the max time the client
	// will sleep before retrying a failed request
	RetryMaxWaitTime = 1 * time.Second

	// Timeout specifies the time a resty client will
	// wait before raising a timeout error.
	Timeout = 500 * time.Millisecond
)

// retryCondFunc is called by resty always and not only on client
// errors. Therefore we need to prevent retrying url errors.
var retryCondFunc = func(r *resty.Response, err error) bool {
	urlErr, ok := err.(*url.Error)

	if ok && !urlErr.Temporary() {
		return false
	}

	if err != nil {
		return true
	}

	return false
}

type HTTPClient interface {
	GET(url string, res interface{}) (*resty.Response, error)
}

type httpClient struct {
	*resty.Client
}

func NewHTTPClient() HTTPClient {
	c := resty.New()
	c.SetRetryCount(RetryCount)
	c.SetRetryWaitTime(RetryWaitTime)
	c.SetRetryMaxWaitTime(RetryMaxWaitTime)
	c.SetTimeout(Timeout)
	c.AddRetryCondition(retryCondFunc)

	return &httpClient{
		Client: c,
	}
}

func (c *httpClient) GET(url string, res interface{}) (*resty.Response, error) {
	req := c.Client.R().SetResult(res)
	httpResp, err := req.Get(url)
	return httpResp, errIfHTTPReqFailed(httpResp, err)
}

func errIfHTTPReqFailed(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("received non 2XX response: %s", resp.Status())
	}

	return nil
}

func buildLatestRateURL(currency string) (string, error) {
	queryParams := map[string]string{
		BaseParam:    currency,
		SymbolsParam: EURSymbol,
	}

	return buildURL(LatestPath, queryParams)
}

func buildHistoricalRatesURL(currency string, startDate, endDate string) (string, error) {
	queryParams := map[string]string{
		StartDateParam: startDate,
		EndDateParam:   endDate,
		SymbolsParam:   EURSymbol,
		BaseParam:      currency,
	}

	return buildURL(HistoricalPath, queryParams)
}

func buildURL(path string, queryParams map[string]string) (string, error) {
	base, err := url.Parse(BaseEndpoint)
	if err != nil {
		return "", err
	}

	// latest path params
	base.Path += path

	// Query params
	params := url.Values{}
	for k, v := range queryParams {
		params.Add(k, v)
	}
	base.RawQuery = params.Encode()

	return base.String(), nil
}
