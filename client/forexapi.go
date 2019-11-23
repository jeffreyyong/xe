package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/jeffreyyong/xe/model"
)

// ForexClient is a client interface for
// calling the https://exchangeratesapi.io/ api.
type ForexClient interface {
	GetLatestRates(rates []string) (*model.LatestRates, error)
	GetHistoricalRates(rates []string, startDate string, endDate string) (*model.HistoricalRates, error)
}

type forexClient struct {
	*resty.Client
}

func NewForexClient() ForexClient {
	restyClient := newDefaultRestyClient()

	return &forexClient{
		Client: restyClient,
	}
}

// GetLatestRates get rates with EUR as base
func (e *forexClient) GetLatestRates(rates []string) (*model.LatestRates, error) {
	url, err := buildLatestRatesURL(rates)
	if err != nil {
		return nil, err
	}

	exchangeRates := &model.LatestRates{}
	req := e.Client.R().SetResult(exchangeRates)
	resp, err := doGet(url, req)
	if err != nil {
		return nil, err
	}

	results, ok := resp.Result().(*model.LatestRates)
	if !ok {
		return nil, NewHTTPClientError(url, "type assertion failed",
			fmt.Errorf("%s: %v", resp.Status(), resp))
	}

	return results, nil
}

// GetHistoricalRates get historical rates with EUR as base
func (e *forexClient) GetHistoricalRates(rates []string, startDate string, endDate string) (*model.HistoricalRates, error) {
	url, err := buildHistoricalRatesURL(rates, startDate, endDate)
	if err != nil {
		return nil, err
	}

	exchangeRates := &model.HistoricalRates{}
	req := e.Client.R().SetResult(exchangeRates)
	resp, err := doGet(url, req)
	if err != nil {
		return nil, err
	}

	results, ok := resp.Result().(*model.HistoricalRates)
	if !ok {
		return nil, NewHTTPClientError(url, "type assertion failed",
			fmt.Errorf("%s: %v", resp.Status(), resp))
	}

	return results, nil
}
