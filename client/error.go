package client

import "fmt"

// HTTPClientError is an error type
// that contains the url, msg and err.
// It implements the golang error interface
// by having the Error() method.
type HTTPClientError struct {
	url string
	msg string
	err error
}

func (e *HTTPClientError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.url, e.msg, e.err)
}

// NewHTTPClientError initialises an HTTPClientError
// given the url, msg and err.
func NewHTTPClientError(url, msg string, err error) error {
	if err == nil {
		return nil
	}

	return &HTTPClientError{url, msg, err}
}
