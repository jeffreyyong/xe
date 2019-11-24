package client

import (
	"errors"

	"github.com/jeffreyyong/xe/model"
)

// Forex is a client interface for
// calling the https://exchangeratesapi.io/ api.
type Forex interface {
	GetLatestRate(currency string) (*model.LatestRate, error)
	GetHistoricalRates(currency string, startDate string, endDate string) (*model.HistoricalRates, error)
}

type forex struct {
	httpClient HTTPClient
}

// NewForex initialises a Forex client
// with a httpClient
func NewForex(c HTTPClient) Forex {
	return &forex{
		httpClient: c,
	}
}

// GetLatestRate gets latest rate from `currency` to EUR
func (e *forex) GetLatestRate(currency string) (*model.LatestRate, error) {
	url, err := buildLatestRateURL(currency)
	if err != nil {
		return nil, err
	}

	rate := &model.LatestRate{}
	resp, err := e.httpClient.GET(url, rate)
	if err != nil {
		return nil, NewHTTPClientError(url, "GetLatestRate", err)
	}

	results, ok := resp.Result().(*model.LatestRate)
	if !ok {
		return nil, NewHTTPClientError(url, "GetLatestRate",
			errors.New("type assertion error"))
	}

	return results, nil
}

// GetHistoricalRates get historical rates from `currency` to EUR
// with the period from the startDate to the endDate
func (e *forex) GetHistoricalRates(currency string, startDate string, endDate string) (*model.HistoricalRates, error) {
	url, err := buildHistoricalRatesURL(currency, startDate, endDate)
	if err != nil {
		return nil, err
	}

	rates := &model.HistoricalRates{}
	resp, err := e.httpClient.GET(url, rates)
	if err != nil {
		return nil, NewHTTPClientError(url, "GetHistoricalRates", err)
	}

	results, ok := resp.Result().(*model.HistoricalRates)
	if !ok {
		return nil, NewHTTPClientError(url, "GetHistoricalRates",
			errors.New("type assertion error"))
	}

	return results, nil
}
