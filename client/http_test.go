package client

import (
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

func TestBuildLatestRatesURL(t *testing.T) {
	rates := []string{"USD", "GBP"}
	expected := "https://api.exchangeratesapi.io/latest?symbols=USD%2CGBP"

	latestRatesURL, err := buildLatestRatesURL(rates)
	assert.NoError(t, err)
	assert.Equal(t, expected, latestRatesURL, "latest rates URL is wrong")
}
