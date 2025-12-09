package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HTTPClient is a mock http client.
type HTTPClient struct {
	mock.Mock
}

// Do performs a mock http request.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
