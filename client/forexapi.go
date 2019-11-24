package client

import (
	"errors"
)

// LatestRate holds the response for the latest rate
// from exchangeratesapi
type LatestRate struct {
	Rates Rates  `json:"rates"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

// HistoricalRates holds the response for the historical rates
// from exchangeratesapi
type HistoricalRates struct {
	RatesList RatesList `json:"rates"`
	Base      string    `json:"base"`
	StartDate string    `json:"start_at"`
	EndDate   string    `json:"end_at"`
}

type Rates map[string]float64

type RatesList map[string]Rates

// Forex is a client interface for
// calling the https://exchangeratesapi.io/ api.
type Forex interface {
	GetLatestRate(currency string) (*LatestRate, error)
	GetHistoricalRates(currency string, startDate string, endDate string) (*HistoricalRates, error)
}

type forex struct {
	httpClient HTTPClient
}

func NewForex(c HTTPClient) Forex {
	return &forex{
		httpClient: c,
	}
}

// GetLatestRate gets latest rate from `currency` to EUR
func (e *forex) GetLatestRate(currency string) (*LatestRate, error) {
	url, err := buildLatestRateURL(currency)
	if err != nil {
		return nil, err
	}

	rate := &LatestRate{}
	resp, err := e.httpClient.GET(url, rate)
	if err != nil {
		return nil, NewHTTPClientError(url, "GetLatestRate", err)
	}

	results, ok := resp.Result().(*LatestRate)
	if !ok {
		return nil, NewHTTPClientError(url, "GetLatestRate",
			errors.New("type assertion error"))
	}

	return results, nil
}

// GetHistoricalRates get historical rates from `currency` to EUR
func (e *forex) GetHistoricalRates(currency string, startDate string, endDate string) (*HistoricalRates, error) {
	url, err := buildHistoricalRatesURL(currency, startDate, endDate)
	if err != nil {
		return nil, err
	}

	rates := &HistoricalRates{}
	resp, err := e.httpClient.GET(url, rates)
	if err != nil {
		return nil, NewHTTPClientError(url, "GetHistoricalRates", err)
	}

	results, ok := resp.Result().(*HistoricalRates)
	if !ok {
		return nil, NewHTTPClientError(url, "GetHistoricalRates",
			errors.New("type assertion error"))
	}

	return results, nil
}
