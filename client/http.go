package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseEndpoint   = "https://api.exchangeratesapi.io"
	PathLatest     = "latest"
	PathHistory    = "history"
	ParamSymbols   = "symbols"
	ParamBase      = "base"
	ParamStartDate = "start_at"
	ParamEndDate   = "end_at"
	SymbolEuro     = "EUR"
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

// HTTPClient is a http client interface
type HTTPClient interface {
	GET(url string, res interface{}) (*resty.Response, error)
}

type httpClient struct {
	*resty.Client
}

// NewHTTPClient returns an instance of resty client
// which implements the HTTPClient interface.
// It also sets some configurations for the client.
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

// GET takes in url and the res interface
// and returns resty Response and error
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

// buildLatestRateURL builds the /latest url given currency
// Note: EUR is always the base so can get the value of 1
// 'currency' in euros
func buildLatestRateURL(currency string) (string, error) {
	queryParams := map[string]string{
		ParamBase:    currency,
		ParamSymbols: SymbolEuro,
	}

	return buildURL(PathLatest, queryParams)
}

// buildHistoricalRatesURL builds the /history url given currency
// startDate and endDate
// Note: EUR is always the base so can get the value of 1
// 'currency' in euros
func buildHistoricalRatesURL(currency string, startDate, endDate string) (string, error) {
	queryParams := map[string]string{
		ParamStartDate: startDate,
		ParamEndDate:   endDate,
		ParamSymbols:   SymbolEuro,
		ParamBase:      currency,
	}

	return buildURL(PathHistory, queryParams)
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
