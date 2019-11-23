package calculator

import (
	"sort"

	"github.com/jeffreyyong/xe/client"
	"gonum.org/v1/gonum/stat"
)

const (
	GBP = "GBP"
	USD = "USD"
)

type Engine interface {
	Inverse(rates client.Rates) client.Rates
}

type engine struct {
}

type Recommendations map[string]bool

type Slopes map[string]float64

func NewEngine() Engine {
	return &engine{}
}

func (c *engine) Inverse(rates client.Rates) client.Rates {
	inverseRates := make(client.Rates, len(rates))
	for k, v := range rates {
		inverseRates[k] = 1 / v
	}
	return inverseRates
}

func sortRatesList(ratesList client.RatesList) []client.Rates {
	var dates []string
	ratesSequence := make([]client.Rates, len(ratesList))

	for d := range ratesList {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	for i, d := range dates {
		ratesSequence[i] = ratesList[d]
	}
	return ratesSequence
}

// linearRegression is
func linearRegression(ratesSequence []client.Rates) Slopes {
	length := len(ratesSequence)
	timeline := make([]float64, length)
	gbpRates := make([]float64, length)
	usdRates := make([]float64, length)

	for i, rates := range ratesSequence {
		timeline[i] = float64(i)
		gbpRates[i] = rates[GBP]
		usdRates[i] = rates[USD]
	}

	_, gbpBeta := stat.LinearRegression(timeline, gbpRates, nil, false)
	_, usdBeta := stat.LinearRegression(timeline, usdRates, nil, false)

	return Slopes{
		GBP: gbpBeta,
		USD: usdBeta,
	}
}
