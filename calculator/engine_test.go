package calculator

import (
	"testing"

	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

// TestSortByDate checks the RatesList is sorted by
// date in ascending order
// Scenario:
// 	- given a rates list with dates in random order
//
// Expect:
// 	- rates are sorted in order in a slice
func TestSortByDate(t *testing.T) {
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

	ratesSequence := sortByDate(ratesList)
	assert.Equal(t, expectedRatesSequence,
		ratesSequence, "rates sequence don't match")
}

// TestGetSlope checks that beta is calculated correctly.
// Scenario:
// 	- given a list of rates that has been ordered
//
// Expect:
// 	- a beta coefficient is calculated correctly
func TestGetSlope(t *testing.T) {
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

// TestRecommend checks that the right signal
// recommendation is given based on the slope
// Scenario:
// 	- explained in the descriptions of tests
//
// Expect:
// 	- right recommendation is given
func TestRecommend(t *testing.T) {
	type testParams struct {
		description       string
		ratesList         model.RatesList
		expRecommendation Signal
	}

	cases := []testParams{
		{
			description: "SignalConvert if price going down",
			ratesList: model.RatesList{
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
			},
			expRecommendation: SignalConvert,
		},
		{
			description: "SignalNoConvert if price is going up",
			ratesList: model.RatesList{
				"2019-11-22": {
					"EUR": 0.1155735337,
				},
				"2019-11-21": {
					"EUR": 0.1152883939,
				},
				"2019-11-20": {
					"EUR": 0.1155414852,
				},
				"2019-11-19": {
					"EUR": 0.115328282,
				},
				"2019-11-18": {
					"EUR": 0.1154827757,
				},
			},
			expRecommendation: SignalNoConvert,
		},
		{
			description: "SignalNeutral if price is constant",
			ratesList: model.RatesList{
				"2019-11-22": {
					"EUR": 0.1155735337,
				},
				"2019-11-21": {
					"EUR": 0.1155735337,
				},
				"2019-11-20": {
					"EUR": 0.1155735337,
				},
				"2019-11-19": {
					"EUR": 0.1155735337,
				},
				"2019-11-18": {
					"EUR": 0.1155735337,
				},
			},
			expRecommendation: SignalNeutral,
		},
	}

	e := NewEngine()
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			recommendation := e.Recommend(tt.ratesList)
			assert.Equal(t, tt.expRecommendation, recommendation,
				"recommendation is wrong")
		})
	}
}
