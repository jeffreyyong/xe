package calculator

import (
	"sort"

	"github.com/jeffreyyong/xe/model"
	"gonum.org/v1/gonum/stat"
)

const (
	EUR = "EUR"

	SignalConvert   Signal = "convert"
	SignalNoConvert Signal = "don't convert"
	SignalNeutral   Signal = "neutral"
)

// Signal is the custom string for signal
type Signal string

// Engine is the calculator interface that
// recommends whether should exchange forex
type Engine interface {
	Recommend(ratesList model.RatesList) Signal
}

type engine struct {
}

// NewEngine initialises the calculator engine
func NewEngine() Engine {
	return &engine{}
}

// Recommend:
// 1) takes a list of rates
// 2) orders them in ascending order by date (exchangeratesapi
//    returns the rates in random order)
// 3) finds the trend line of those rates by calculating
//    the beta
// 4) returns 'convert' if the price is cheaper, returns 'don't
//    convert' if the price is more expensive, and 'neutral'
//    if the price is constant.
func (e *engine) Recommend(ratesList model.RatesList) Signal {
	sortedRates := sortByDate(ratesList)
	slope := getSlope(sortedRates, EUR)

	signal := SignalNeutral
	if slope > 0 {
		signal = SignalNoConvert
	} else if slope < 0 {
		signal = SignalConvert
	}

	return signal
}

// sortByDate sorts the rates by the key of RatesList, i.e. date string
func sortByDate(ratesList model.RatesList) []model.Rates {
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

// getSlope computes the beta (trend line) given a list of ordered rates.
func getSlope(ratesSequence []model.Rates, currency string) float64 {
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
