package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/jeffreyyong/xe/client/mock"
	clientmock "github.com/jeffreyyong/xe/client/mock"
	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

// TestGetLatestRateHappyCase tests that the latest rate
// is returned when given a currency
// Scenario:
// 	- the httpClient is mocked to return a LatestRate object
// 	- and no error is returned
//
// Expect:
// 	- asserts that no error is return by GetLatestRate
// 	- asserts that the right LatestRate result is returned
func TestGetLatestRateHappyCase(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

	mockLatestRate := &model.LatestRate{
		Rates: model.Rates{
			"EUR": 1.163061177,
		},
		Base: "GBP",
		Date: "2019-11-22",
	}
	mockHTTPClientResp := &resty.Response{
		Request: &resty.Request{
			Result: mockLatestRate,
		},
	}
	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(mockHTTPClientResp, nil)

	latestRate, err := forex.GetLatestRate("GBP")
	assert.NoError(t, err)
	assert.Equal(t, mockLatestRate, latestRate, "result does not match")
}

// TestGetLatestRateHTTPClientError tests that an error is returned
// by the method when httpClient has error
// Scenario:
// 	- error returned by the httpClient
//
// Expect:
// 	- right HTTPClientError is returned by the method GetLatestRate
//  - result is nil
func TestGetLatestRateHTTPClientError(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

	errString := "connection closed"
	mockHTTPClientErr := errors.New(errString)
	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(nil, mockHTTPClientErr)
	latestRate, err := forex.GetLatestRate("GBP")

	expectedErrString := "https://api.exchangeratesapi.io/latest?base=GBP&symbols=EUR: GetLatestRate: connection closed"
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErrString)
	assert.Nil(t, latestRate)
}

// TestGetLatestRateTypeAssertion tests that type assertion error
// is raised by the method when the result can't be asserted to type LatestRate
// Scenario:
// 	- httpClient returns a result interface that is string "foobar" which is not the
// 	  the type of LatestRate
//
// Expect:
// 	- type assertion error is returned in the format of HTTPClientError
func TestGetLatestRateTypeAssertion(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

	mockHTTPClientResp := &resty.Response{
		Request: &resty.Request{
			Result: "foobar",
		},
		RawResponse: &http.Response{
			Status: "200 OK",
		},
	}

	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(mockHTTPClientResp, nil)
	latestRate, err := forex.GetLatestRate("GBP")

	expectedErrString := "https://api.exchangeratesapi.io/latest?base=GBP&symbols=EUR: GetLatestRate: type assertion error"
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErrString)
	assert.Nil(t, latestRate)
}

// TestGetHistoricalRatesHappyCase tests that historal rates
// are returned when given the currency, startDate and endDate
// Scenario:
//	- the httpClient is mocked to return a HistoricalRates object
// 	- and no error is returned
//
// Expect:
// 	- asserts that no error is returned by GetHistoricalRates
// 	- asserts that the right HistoricalRates result is returned
func TestGetHistoricalRatesHappyCase(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

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

	mockHTTPClientResp := &resty.Response{
		Request: &resty.Request{
			Result: mockHistoricalRates,
		},
	}

	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(mockHTTPClientResp, nil)

	rates, err := forex.GetHistoricalRates("GBP", "2019-11-21", "2019-11-22")
	assert.NoError(t, err)
	assert.Equal(t, mockHistoricalRates, rates, "result does not match")
}

// TestGetHistoricalRatesHTTPClientError tests that an error is returned
// by the method when httpClient has error
// Scenario:
// 	- error returned by the httpClient
//
// Expect:
// 	- right HTTPClientError is returned by the method GetHistoricalRates
//  - result is nil
func TestGetHistoricalRatesHTTPClientError(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

	errString := "connection closed"
	mockHTTPClientErr := errors.New(errString)
	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(nil, mockHTTPClientErr)
	latestRate, err := forex.GetHistoricalRates("GBP", "2019-11-21", "2019-11-22")

	expectedErrString := "https://api.exchangeratesapi.io/history?base=GBP&end_at=2019-11-22&start_at=2019-11-21&symbols=EUR: GetHistoricalRates: connection closed"
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErrString)
	assert.Nil(t, latestRate)
}

// TestGetHistoricalRatesTypeAssertion tests that type assertion error
// is raised by the method when the result object fails to be asserted to type HistoricalRates
// Scenario:
// 	- httpClient returns a result interface that is string "foobar" which is not the
// 	  the type of HistoricalRates
//
// Expect:
// 	- type assertion error is returned in the format of HTTPClientError
func TestGetHistoricalRatesTypeAssertion(t *testing.T) {
	httpClient, forex, ctrl := setupTestForex(t)
	defer ctrl.Finish()

	mockHTTPClientResp := &resty.Response{
		Request: &resty.Request{
			Result: "foobar",
		},
		RawResponse: &http.Response{
			Status: "200 OK",
		},
	}

	httpClient.EXPECT().GET(gomock.Any(), gomock.Any()).
		Return(mockHTTPClientResp, nil)
	latestRate, err := forex.GetHistoricalRates("GBP", "2019-11-21", "2019-11-22")

	expectedErrString := "https://api.exchangeratesapi.io/history?base=GBP&end_at=2019-11-22&start_at=2019-11-21&symbols=EUR: GetHistoricalRates: type assertion error"
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErrString)
	assert.Nil(t, latestRate)
}

func setupTestForex(t *testing.T) (*clientmock.MockHTTPClient, Forex, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)
	return httpClient, forex, ctrl
}
