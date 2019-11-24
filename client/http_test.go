package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// TestRestyClientRetryCondition checks the client
// does not retry on malformed URLs
// Scenario:
// 	- wrap the default retry condition function,
//	- set the retry count to 10
//
// Expect:
//	- the wrapper is called exactly once
func TestRestyClientRetryCondition(t *testing.T) {
	gotRetries := 0
	expRetries := 1

	mockRetryCondFunc := func(r *resty.Response, err error) bool {
		gotRetries++
		assert.Equal(t, expRetries, gotRetries)

		return retryCondFunc(r, err)
	}

	client := resty.New()
	client.SetRetryCount(10)
	client.AddRetryCondition(mockRetryCondFunc)

	_, err := client.R().Get("//")
	assert.Error(t, err)
}

// TestRestyClientRetryCondition checks the client
// does not retry on 200 OK
//
// Scenario:
//   wrap the default retry condition function
//
// Expect:
//   the wrapper is called exactly once
func TestRestyClientRetryConditionNoRetries(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(http.Handler(handler))
	defer ts.Close()

	gotRetries := 0
	expRetries := 1

	mockRetryCondFunc := func(r *resty.Response, err error) bool {
		gotRetries++
		assert.Equal(t, expRetries, gotRetries)

		return retryCondFunc(r, err)
	}

	client := resty.New()
	client.SetRetryCount(2)
	client.AddRetryCondition(mockRetryCondFunc)

	_, err := client.R().Get(ts.URL)
	assert.NoError(t, err)
}

func TestBuildLatestRateURL(t *testing.T) {
	rate := "GBP"
	expected := "https://api.exchangeratesapi.io/latest?base=GBP&symbols=EUR"

	latestRateURL, err := buildLatestRateURL(rate)
	assert.NoError(t, err)
	assert.Equal(t, expected, latestRateURL, "latest rate URL is wrong")
}

func TestBuildHistoricalRatesURL(t *testing.T) {
	rate := "GBP"
	startDate := "2019-11-15"
	endDate := "2019-11-22"
	expected := "https://api.exchangeratesapi.io/history?base=GBP&end_at=2019-11-22&start_at=2019-11-15&symbols=EUR"

	latestRateURL, err := buildHistoricalRatesURL(rate, startDate, endDate)
	assert.NoError(t, err)
	assert.Equal(t, expected, latestRateURL, "historical rates URL is wrong")
}

func TestErrIfHTTPReqFailed(t *testing.T) {
	type testParams struct {
		description   string
		err           error
		resp          *resty.Response
		expErrExist   bool
		expErrMessage string
	}

	cases := []testParams{
		{
			description:   "err exists before HTTP request is made",
			err:           errors.New("tcp connection timeout"),
			resp:          nil,
			expErrExist:   true,
			expErrMessage: "tcp connection timeout",
		},
		{
			description: "non 2xx status code",
			err:         nil,
			resp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: 500,
					Status:     "500 ERROR",
				},
			},
			expErrExist:   true,
			expErrMessage: "received non 2XX response: 500 ERROR",
		},
		{
			description: "no error",
			err:         nil,
			resp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: 200,
					Status:     "200 OK",
				},
			},
			expErrExist: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			err := errIfHTTPReqFailed(tt.resp, tt.err)
			if tt.expErrExist {
				assert.Error(t, err)
				assert.Equal(t, tt.expErrMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
