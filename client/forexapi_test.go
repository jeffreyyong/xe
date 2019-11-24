package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/jeffreyyong/xe/client/mock"
	"github.com/jeffreyyong/xe/model"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestRateHappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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

func TestGetLatestRateHTTPClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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

func TestGetLatestRateTypeAssertion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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

func TestGetHistoricalRatesHappyCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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

func TestGetHistoricalRatesHTTPClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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

func TestGetHistoricalRatesTypeAssertion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := mock.NewMockHTTPClient(ctrl)
	forex := NewForex(httpClient)

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
