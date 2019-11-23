package calculator

import (
	"sort"

	"github.com/jeffreyyong/xe/model"
	"gonum.org/v1/gonum/stat"
)

const (
	GBP = "GBP"
	USD = "USD"
)

type Calculator interface {
	Inverse(rates model.Rates) model.Rates
}

type calculator struct {
}

type Recommendations map[string]bool

type Slopes map[string]float64

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) Inverse(rates model.Rates) model.Rates {
	inverseRates := make(model.Rates, len(rates))
	for k, v := range rates {
		inverseRates[k] = 1 / v
	}
	return inverseRates
}

func sortRatesList(ratesList model.RatesList) []model.Rates {
	var dates []string
	ratesSequence := make([]model.Rates, len(ratesList))

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
func linearRegression(ratesSequence []model.Rates) Slopes {
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
