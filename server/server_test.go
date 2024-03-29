package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jeffreyyong/xe/calculator"
	calculatormock "github.com/jeffreyyong/xe/calculator/mock"
	"github.com/jeffreyyong/xe/client"
	forexmock "github.com/jeffreyyong/xe/client/mock"
	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

const (
	testServerAddr = "localhost:3000"
)

// TestQueryParamsMissing validates against missing query params
// Scenario:
// 	- query params of 'currency' not provided
//
// Expect:
// 	- right error message is returned in the JSON body
// 	- StatusCode of 400 is returned
func TestQueryParamsMissing(t *testing.T) {
	_, _, xeService, ctrl := setupTestServer(t)
	defer ctrl.Finish()

	go xeService.Run()
	defer xeService.Stop()

	convertResp := &model.ConvertResp{}
	urlNoQueryParam := "http://localhost:3000/convert"
	httpClient := client.NewHTTPClient()
	resp, err := httpClient.GET(urlNoQueryParam, convertResp)

	expJSON := `{"error":"invalid query parameter - currency must be provided"}`
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	assert.Equal(t, expJSON, string(resp.Body()))
}

// TestGetLatestRateError checks if error is returned
// when fx.GetLatestRate returns error
// Scenario:
// 	- mockFX is configured to return an error for GetLatestRate
//
// Expect:
// 	- right error message is returned in the JSON body
// 	- StatusCode of 500 is returned
func TestGetLatestRateError(t *testing.T) {
	_, mockFX, xeService, ctrl := setupTestServer(t)
	defer ctrl.Finish()

	go xeService.Run()
	defer xeService.Stop()

	mockFX.EXPECT().GetLatestRate(gomock.Any()).
		Return(nil, errors.New("error getting latest rate"))

	convertResp := &model.ConvertResp{}
	url := "http://localhost:3000/convert?currency=USD"
	httpClient := client.NewHTTPClient()
	resp, err := httpClient.GET(url, convertResp)

	expJSON := `{"error":"error converting currency"}`
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	assert.Equal(t, expJSON, string(resp.Body()))
}

// TestExtractTargetRateError checks if error is returned
// when target rate can't be retrieved from LatestRate
// Scenario:
// 	- model.Rates contains an unknown key, i.e. not EUR
//
// Expect:
// 	- right error message is returned in the JSON body
// 	- StatusCode of 500 is returned
func TestExtractTargetRateError(t *testing.T) {
	_, mockFX, xeService, ctrl := setupTestServer(t)
	defer ctrl.Finish()

	go xeService.Run()
	defer xeService.Stop()

	// error with unrecognised interest rate:
	// NON_EXISTENT_RATE
	mockLatestRate := &model.LatestRate{
		Rates: model.Rates{
			"NON_EXISTENT_RATE": 1.163061177,
		},
		Base: "GBP",
		Date: "2019-11-22",
	}

	mockFX.EXPECT().GetLatestRate(gomock.Any()).
		Return(mockLatestRate, nil)

	convertResp := &model.ConvertResp{}
	url := "http://localhost:3000/convert?currency=USD"
	httpClient := client.NewHTTPClient()
	resp, err := httpClient.GET(url, convertResp)

	expJSON := `{"error":"error converting currency"}`
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	assert.Equal(t, expJSON, string(resp.Body()))
}

// TestGetHistoricalRatesError checks if eror is returned
// when fx.GetHistoricalRates returns error
// Scenario:
// 	- mockFX is configured to return error for GetHistoricalRates
//
// Expect:
// 	- right error message is returned in the JSON body
// 	- StatusCode of 500 is returned
func TestGetHistoricalRatesError(t *testing.T) {
	_, mockFX, xeService, ctrl := setupTestServer(t)
	defer ctrl.Finish()

	go xeService.Run()
	defer xeService.Stop()

	mockLatestRate := &model.LatestRate{
		Rates: model.Rates{
			"EUR": 1.163061177,
		},
		Base: "GBP",
		Date: "2019-11-22",
	}

	mockFX.EXPECT().GetLatestRate(gomock.Any()).
		Return(mockLatestRate, nil)

	mockFX.EXPECT().GetHistoricalRates(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("error getting historical rate"))

	convertResp := &model.ConvertResp{}
	url := "http://localhost:3000/convert?currency=USD"
	httpClient := client.NewHTTPClient()
	resp, err := httpClient.GET(url, convertResp)

	expJSON := `{"error":"error converting currency"}`
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	assert.Equal(t, expJSON, string(resp.Body()))
}

// TestHandlerConvertNoError shows the happy path
// Scenario:
// 	- mockFX returns correct mockLatestRate
// 	- mockFX returns correct mockHistoricalRates
//  - mockCE gives recommendation
//
// Expect:
// 	- right JSON is provided: {"from":"USD","to":"EUR","rate":1.163061177,"recommendation":"convert"}
// 	- StatusCode of 200 is returned
func TestHandlerConvertNoError(t *testing.T) {
	mockCE, mockFX, xeService, ctrl := setupTestServer(t)
	defer ctrl.Finish()

	go xeService.Run()
	defer xeService.Stop()

	mockLatestRate := &model.LatestRate{
		Rates: model.Rates{
			"EUR": 1.163061177,
		},
		Base: "GBP",
		Date: "2019-11-22",
	}

	mockHistoricalRates := &model.HistoricalRates{
		RatesList: model.RatesList{
			"2019-11-21": model.Rates{
				"EUR": 1.163061177,
			},
			"2019-11-22": model.Rates{
				"EUR": 1.163061177,
			},
		},
		Base:      "GBP",
		StartDate: "2019-11-21",
		EndDate:   "2019-11-22",
	}

	mockFX.EXPECT().GetLatestRate(gomock.Any()).
		Return(mockLatestRate, nil)

	mockFX.EXPECT().GetHistoricalRates(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockHistoricalRates, nil)

	mockCE.EXPECT().Recommend(gomock.Any()).Return(calculator.SignalConvert)

	convertResp := &model.ConvertResp{}
	url := "http://localhost:3000/convert?currency=USD"
	httpClient := client.NewHTTPClient()
	resp, err := httpClient.GET(url, convertResp)

	expJSON := `{"from":"USD","to":"EUR","rate":1.163061177,"recommendation":"convert"}`
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, expJSON, string(resp.Body()))

}

func setupTestServer(t *testing.T) (*calculatormock.MockEngine, *forexmock.MockForex, *XEService, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockFX := forexmock.NewMockForex(ctrl)
	mockCE := calculatormock.NewMockEngine(ctrl)

	h := NewHandler(mockFX, mockCE)
	httpHandler := SetupAPIHandler(h)
	xeService := NewXEService(httpHandler, testServerAddr)

	return mockCE, mockFX, xeService, ctrl
}
