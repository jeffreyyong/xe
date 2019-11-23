package client

import "fmt"

type HTTPClientError struct {
	url string
	msg string
	err error
}

func (e *HTTPClientError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.url, e.msg, e.err)
}

func NewHTTPClientError(url, msg string, err error) error {
	if err == nil {
		return nil
	}

	return &HTTPClientError{url, msg, err}
}
