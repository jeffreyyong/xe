package model

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
