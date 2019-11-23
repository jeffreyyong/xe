package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParams struct {
	description         string
	url                 string
	msg                 string
	err                 error
	expHTTPClientErrNil bool
	expErrMsg           string
}

func TestNewHTTPClientError(t *testing.T) {
	cases := []testParams{
		{
			description:         "non nil error",
			url:                 "http://localhost.com",
			msg:                 "GetLatestRates",
			err:                 errors.New("interface conversion: interface is string not int"),
			expHTTPClientErrNil: false,
			expErrMsg:           "http://localhost.com: GetLatestRates: interface conversion: interface is string not int",
		},
		{
			description:         "nil error",
			url:                 "http://localhost.com",
			msg:                 "",
			err:                 nil,
			expHTTPClientErrNil: true,
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
	if tt.expHTTPClientErrNil {
		assert.Nil(t, httpClientErr)
	} else {
		assert.NotNil(t, httpClientErr)
		assert.Equal(t, tt.expErrMsg, httpClientErr.Error(), "error message is wrong")
	}
}
