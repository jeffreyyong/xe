package model

// LatestRate holds the response for the latest rate
// from exchangeratesapi
type LatestRate struct {
	Rates Rates  `json:"rates"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

// HistoricalRates holds the response for the history rates
// from exchangeratesapi
type HistoricalRates struct {
	RatesList RatesList `json:"rates"`
	Base      string    `json:"base"`
	StartDate string    `json:"start_at"`
	EndDate   string    `json:"end_at"`
}

// Rates is a map of currency:rate
// e.g.
// {
//   "EUR": 1.163061177
// },
type Rates map[string]float64

// Rateslist is a map of date:rates
// e.g.
// {
// 	 "2019-11-21": {
// 	 	"EUR": 1.1689343994
// 	 },
// 	 "2019-11-15": {
// 	 	"EUR": 1.1674060238
// 	 },
// 	 "2019-11-22": {
// 	 	"EUR": 1.163061177
// 	 },
// 	 "2019-11-20": {
// 	 	"EUR": 1.1666569445
// 	 },
// 	 "2019-11-19": {
// 	 	"EUR": 1.1685928973
// 	 },
// 	 "2019-11-18": {
// 	 	"EUR": 1.1719207782
// 	 }
// },
type RatesList map[string]Rates
