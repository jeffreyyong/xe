package calculator

import (
	"sort"

	"github.com/jeffreyyong/xe/client"
	"gonum.org/v1/gonum/stat"
)

const (
	EUR = "EUR"

	ConvertSignal   Signal = "convert"
	NoConvertSignal Signal = "don't convert"
	NeutralSignal   Signal = "Neutral"
)

type Signal string

type Engine interface {
	Recommend(ratesList client.RatesList) Signal
}

type engine struct {
}

func NewEngine() Engine {
	return &engine{}
}

func (e *engine) Recommend(ratesList client.RatesList) Signal {
	sortedRates := sortRatesList(ratesList)
	slope := getSlope(sortedRates, EUR)

	signal := NeutralSignal
	if slope > 0 {
		signal = NoConvertSignal
	} else if slope < 0 {
		signal = ConvertSignal
	}

	return signal
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

// getSlope gets the trend of a line
func getSlope(ratesSequence []client.Rates, currency string) float64 {
	length := len(ratesSequence)
	timeline := make([]float64, length)
	rates := make([]float64, length)

	for i, r := range ratesSequence {
		timeline[i] = float64(i)
		rates[i] = r[currency]

	}

	_, beta := stat.LinearRegression(timeline, rates, nil, false)
	return beta
}