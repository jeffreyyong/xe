package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// LatestRates holds the response for the latest rates
type LatestRates struct {
	Rates Rates  `json:"rates"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

// HistoricalRates holds the response for the historical rates
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
	GetLatestRates(rates []string) (*LatestRates, error)
	GetHistoricalRates(rates []string, startDate string, endDate string) (*HistoricalRates, error)
}

type forex struct {
	*resty.Client
}

func NewForexClient() Forex {
	restyClient := newDefaultRestyClient()

	return &forex{
		Client: restyClient,
	}
}

// GetLatestRates get rates with EUR as base
func (e *forex) GetLatestRates(rates []string) (*LatestRates, error) {
	url, err := buildLatestRatesURL(rates)
	if err != nil {
		return nil, err
	}

	exchangeRates := &LatestRates{}
	req := e.Client.R().SetResult(exchangeRates)
	resp, err := doGet(url, req)
	if err != nil {
		return nil, err
	}

	results, ok := resp.Result().(*LatestRates)
	if !ok {
		return nil, NewHTTPClientError(url, "type assertion failed",
			fmt.Errorf("%s: %v", resp.Status(), resp))
	}

	return results, nil
}

// GetHistoricalRates get historical rates with EUR as base
func (e *forex) GetHistoricalRates(rates []string, startDate string, endDate string) (*HistoricalRates, error) {
	url, err := buildHistoricalRatesURL(rates, startDate, endDate)
	if err != nil {
		return nil, err
	}

	exchangeRates := &HistoricalRates{}
	req := e.Client.R().SetResult(exchangeRates)
	resp, err := doGet(url, req)
	if err != nil {
		return nil, err
	}

	results, ok := resp.Result().(*HistoricalRates)
	if !ok {
		return nil, NewHTTPClientError(url, "type assertion failed",
			fmt.Errorf("%s: %v", resp.Status(), resp))
	}

	return results, nil
}
