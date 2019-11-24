package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParams struct {
	description            string
	url                    string
	msg                    string
	err                    error
	expHTTPClientErrExists bool
	expErrMsg              string
}

// TestNewHTTPClientError checks that the right format
// of error is constructed
func TestNewHTTPClientError(t *testing.T) {
	cases := []testParams{
		{
			description:            "non nil error",
			url:                    "http://localhost.com",
			msg:                    "GetLatestRate",
			err:                    errors.New("interface conversion: interface is string not int"),
			expHTTPClientErrExists: true,
			expErrMsg:              "http://localhost.com: GetLatestRate: interface conversion: interface is string not int",
		},
		{
			description:            "nil error",
			url:                    "http://localhost.com",
			msg:                    "",
			err:                    nil,
			expHTTPClientErrExists: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			httpClientErr := NewHTTPClientError(tt.url, tt.msg, tt.err)
			assertTest(t, httpClientErr, tt)
		})
	}
}

func assertTest(t *testing.T, httpClientErr error, tt testParams) {
	if tt.expHTTPClientErrExists {
		assert.Error(t, httpClientErr)
		assert.Equal(t, tt.expErrMsg, httpClientErr.Error(), "error message is wrong")
	} else {
		assert.NoError(t, httpClientErr)
	}
}
