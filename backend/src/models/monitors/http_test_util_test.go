package monitors

import (
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of the httpClient interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// createMockResponse creates a mock http.Response for testing
func createMockResponse(statusCode int, body string, headers map[string]string) *http.Response {
	header := make(http.Header)
	for k, v := range headers {
		header.Add(k, v)
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
	}
}

// setupTestHttpMonitorConfig returns a default HttpConfig for testing
func setupTestHttpMonitorConfig() *HttpConfig {
	return &HttpConfig{
		Method:              "GET",
		URL:                 "http://example.com",
		ExpectedStatusCodes: []int{200},
	}
}
