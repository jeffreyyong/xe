package calculator

import (
	"testing"

	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

func TestSortRatesList(t *testing.T) {
	ratesList := model.RatesList{
		"2019-11-21": {
			EUR: 1.1689343994,
		},
		"2019-11-15": {
			EUR: 1.1674060238,
		},
		"2019-11-22": {
			EUR: 1.163061177,
		},
		"2019-11-20": {
			EUR: 1.1666569445,
		},
		"2019-11-19": {
			EUR: 1.1685928973,
		},
		"2019-11-18": {
			EUR: 1.1719207782,
		},
	}

	expectedRatesSequence := []model.Rates{
		{
			EUR: 1.1674060238,
		},
		{
			EUR: 1.1719207782,
		},
		{
			EUR: 1.1685928973,
		},
		{
			EUR: 1.1666569445,
		},
		{
			EUR: 1.1689343994,
		},
		{
			EUR: 1.163061177,
		},
	}

	ratesSequence := sortRatesList(ratesList)
	assert.Equal(t, expectedRatesSequence,
		ratesSequence, "rates sequence don't match")
}

func TestLinearRegression(t *testing.T) {
	ratesSequence := []model.Rates{
		{
			EUR: 1.1674060238,
		},
		{
			EUR: 1.1719207782,
		},
		{
			EUR: 1.1685928973,
		},
		{
			EUR: 1.1666569445,
		},
		{
			EUR: 1.1689343994,
		},
		{
			EUR: 1.163061177,
		},
	}

	expectedSlope := -0.0009319806628571443

	slope := getSlope(ratesSequence, EUR)
	assert.Equal(t, expectedSlope,
		slope, "slope is wrong")
}
