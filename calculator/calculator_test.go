package calculator

import (
	"testing"

	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

func TestInverse(t *testing.T) {
	c := NewCalculator()
	rates := model.Rates{
		USD: 1.1058,
		GBP: 0.8598,
	}
	expectedInverseRates := model.Rates{
		USD: 0.904322662325918,
		GBP: 1.1630611770179111,
	}

	inverseRates := c.Inverse(rates)
	assert.Equal(t, expectedInverseRates,
		inverseRates, "inverse rates are wrong")
}

func TestSortRatesList(t *testing.T) {
	ratesList := model.RatesList{
		"2019-11-21": {
			USD: 1.1091,
			GBP: 0.85548,
		},
		"2019-11-15": {
			USD: 1.1034,
			GBP: 0.8566,
		},
		"2019-11-22": {
			USD: 1.1058,
			GBP: 0.8598,
		},
		"2019-11-20": {
			USD: 1.1059,
			GBP: 0.85715,
		},
		"2019-11-19": {
			USD: 1.1077,
			GBP: 0.85573,
		},
		"2019-11-18": {
			USD: 1.1061,
			GBP: 0.8533,
		},
	}

	expectedRatesSequence := []model.Rates{
		{
			USD: 1.1034,
			GBP: 0.8566,
		},
		{
			USD: 1.1061,
			GBP: 0.8533,
		},
		{
			USD: 1.1077,
			GBP: 0.85573,
		},
		{
			USD: 1.1059,
			GBP: 0.85715,
		},
		{
			USD: 1.1091,
			GBP: 0.85548,
		},
		{
			USD: 1.1058,
			GBP: 0.8598,
		},
	}

	ratesSequence := sortRatesList(ratesList)
	assert.Equal(t, expectedRatesSequence,
		ratesSequence, "rates sequence don't match")
}

func TestLinearRegression(t *testing.T) {
	ratesSequence := []model.Rates{
		{
			USD: 1.1034,
			GBP: 0.8566,
		},
		{
			USD: 1.1061,
			GBP: 0.8533,
		},
		{
			USD: 1.1077,
			GBP: 0.85573,
		},
		{
			USD: 1.1059,
			GBP: 0.85715,
		},
		{
			USD: 1.1091,
			GBP: 0.85548,
		},
		{
			USD: 1.1058,
			GBP: 0.8598,
		},
	}

	expectedSlopes := Slopes{
		GBP: 0.0006845714285714312,
		USD: 0.0005485714285714189,
	}

	slopes := linearRegression(ratesSequence)
	assert.Equal(t, expectedSlopes,
		slopes, "slopes are wrong")
}
