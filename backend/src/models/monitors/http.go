package monitors

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/models/shared"
	"github.com/m-milek/leszmonitor/util"
)

type HTTPProbe struct {
	Method               string            `json:"method"`
	URL                  string            `json:"url"`
	Headers              map[string]string `json:"headers"`
	Body                 string            `json:"body"`
	SaveResponseBody     bool              `json:"saveResponseBody"`    // Whether to save the response body in the monitor response
	SaveResponseHeaders  bool              `json:"saveResponseHeaders"` // Whether to save the response headers in the monitor response
	ExpectedStatusCodes  []int             `json:"expectedStatusCodes"`
	ExpectedBodyRegex    string            `json:"expectedBodyRegex"`
	ExpectedHeaders      map[string]string `json:"expectedHeaders"`
	ExpectedResponseTime *int              `json:"expectedResponseTime"` // in milliseconds
}

const httpTimeout = 10 * time.Second

func (m *HTTPProbe) Run(ctx context.Context, monitorID uuid.UUID) monitorresult.IMonitorResult {
	logger := log.FromContext(ctx)
	result := monitorresult.NewMonitorResult(
		monitorID,
		consts.HTTPConfigType,
		shared.MonitorStatusUp,
		false,
		0,
		"",
		&monitorresult.HTTPResultDetails{},
	)
	details, castErr := result.GetDetails().(*monitorresult.HTTPResultDetails)
	if !castErr {
		logger.Error().Msg("Failed to cast monitor result details to HTTPResultDetails")
		result.AddError("Internal error: failed to process monitor result details")
		return &result
	}

	httpResponse, elapsed, err := m.executeRequest(&httpClientOrMock)
	if httpResponse != nil {
		defer httpResponse.Body.Close()
	}

	result.SetDuration(elapsed.Milliseconds())
	if err != nil {
		result.AddError(fmt.Sprintf("HTTP request failed: %s", err.Error()))
		logger.Trace().Err(err).Msg("HTTP request execution failed")
		return &result
	}

	details.StatusCode = httpResponse.StatusCode
	details.Proto = httpResponse.Proto
	details.ContentLength = httpResponse.ContentLength

	if m.SaveResponseHeaders {
		details.Headers = make(map[string]string)
		for key, value := range httpResponse.Header {
			details.Headers[key] = strings.Join(value, ", ")
		}
	}

	if m.SaveResponseBody {
		body, err := readResponseBody(httpResponse)
		if err == nil {
			details.Body = body
		}
	}

	m.checkStatusCode(httpResponse, &result)
	m.checkResponseTime(elapsed, &result)
	m.checkResponseHeaders(httpResponse, &result)
	m.checkResponseBody(httpResponse, &result)

	return &result
}

// Validate checks if the HTTP monitor configuration is valid
// It ensures that required fields are set and that the URL is properly formatted.
func (m *HTTPProbe) Validate() error {
	if m.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if m.Method == "" {
		return fmt.Errorf("HTTP method cannot be empty")
	}

	if !util.SliceContains([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}, m.Method) {
		return fmt.Errorf("invalid HTTP method: %s", m.Method)
	}

	if len(m.ExpectedStatusCodes) == 0 {
		return fmt.Errorf("expected status codes cannot be empty")
	}

	if len(m.ExpectedStatusCodes) > 0 {
		minValue, maxValue := util.SliceMinMax(m.ExpectedStatusCodes)
		if minValue < 100 || maxValue > 599 {
			return fmt.Errorf("expected status codes must be between 100 and 599")
		}
	}

	parsedURL, err := url.Parse(m.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid URL format: %s", m.URL)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be either http or https: %s", m.URL)
	}

	if m.ExpectedResponseTime != nil && *m.ExpectedResponseTime < 0 {
		return fmt.Errorf("expected response time cannot be negative")
	}

	if m.ExpectedBodyRegex != "" {
		if _, err := regexp.Compile(m.ExpectedBodyRegex); err != nil {
			return fmt.Errorf("invalid body regex: %w", err)
		}
	}

	return nil
}

// httpRequestExecutor is needed for mocking HTTP requests in tests.
type httpRequestExecutor interface {
	Do(req *http.Request) (*http.Response, error)
}

func newHTTPClient() httpRequestExecutor {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

var httpClientOrMock = newHTTPClient()

// Encapsulates request creation and execution.
func (m *HTTPProbe) executeRequest(httpClient *httpRequestExecutor) (*http.Response, time.Duration, error) {
	request, err := m.createRequest()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("User-Agent", "LeszMonitor/DEV")

	start := time.Now()
	response, err := (*httpClient).Do(request)
	elapsed := time.Since(start)

	if err != nil {
		return nil, elapsed, err
	}

	return response, elapsed, nil
}

func (m *HTTPProbe) checkStatusCode(
	response *http.Response,
	result monitorresult.IMonitorResult,
) {
	if m.ExpectedStatusCodes == nil {
		return
	}

	if !util.SliceContains(m.ExpectedStatusCodes, response.StatusCode) {
		failureMsg := fmt.Sprintf(
			"Unexpected status code: got %d, expected one of %v",
			response.StatusCode,
			m.ExpectedStatusCodes,
		)
		result.AddFailure(failureMsg)
	}
}

func (m *HTTPProbe) checkResponseTime(
	elapsed time.Duration,
	result monitorresult.IMonitorResult,
) {
	if m.ExpectedResponseTime == nil {
		return
	}
	if elapsed.Milliseconds() > int64(*m.ExpectedResponseTime) {
		failureMsg := fmt.Sprintf(
			"Response time exceeded: got %dms, expected <= %dms",
			elapsed.Milliseconds(),
			*m.ExpectedResponseTime,
		)
		result.AddFailure(failureMsg)
	}
}

func (m *HTTPProbe) checkResponseHeaders(
	response *http.Response,
	result monitorresult.IMonitorResult,
) {
	if len(m.ExpectedHeaders) == 0 {
		return
	}

	for key, expectedValue := range m.ExpectedHeaders {
		actualValue := response.Header.Get(key)
		if actualValue != expectedValue {
			failureMsg := fmt.Sprintf("Header mismatch for %s: got %s, expected %s", key, actualValue, expectedValue)
			result.AddFailure(failureMsg)
		}
	}
}

func (m *HTTPProbe) checkResponseBody(
	response *http.Response,
	result monitorresult.IMonitorResult,
) {
	if m.ExpectedBodyRegex == "" {
		return
	}

	responseBody, err := readResponseBody(response)
	if err != nil {
		result.AddError("Error reading response body: " + err.Error())
		return
	}

	// Add (?s) flag to make dot match newlines
	patternWithFlag := "(?s)" + m.ExpectedBodyRegex

	regex, err := regexp.Compile(patternWithFlag)
	if err != nil {
		result.AddError(fmt.Sprintf("Invalid regex for expected body: %s", patternWithFlag))
		return
	}

	matches := regex.MatchString(responseBody)
	if !matches {
		failureMsg := fmt.Sprintf("Response body does not match regex: %s", m.ExpectedBodyRegex)
		result.AddFailure(failureMsg)
	}
}

// createRequest constructs an HTTP request based on the monitor's configuration.
func (m *HTTPProbe) createRequest() (*http.Request, error) {
	parsedURL, err := url.Parse(m.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", m.URL)
	}

	req := http.Request{
		Method: m.Method,
		URL:    parsedURL,
		Header: make(http.Header),
	}

	for key, value := range m.Headers {
		req.Header.Set(key, value)
	}

	if m.Body != "" {
		req.Body = io.NopCloser(strings.NewReader(m.Body))
	} else {
		req.Body = nil
	}

	return &req, nil
}

// Helper function to read response body while preserving it.
func readResponseBody(response *http.Response) (string, error) {
	if response.Body == nil {
		return "", nil
	}

	// Read the body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Restore the body so it can be read again
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}
