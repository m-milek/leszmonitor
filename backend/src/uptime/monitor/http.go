package monitors

import (
	"bytes"
	"fmt"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type httpConfig struct {
	Method               string            `json:"httpMethod" bson:"httpMethod"`
	URL                  string            `json:"url" bson:"url"`
	Headers              map[string]string `json:"headers" bson:"headers"`
	Body                 string            `json:"body" bson:"body"`
	SaveResponseBody     bool              `json:"saveResponseBody" bson:"saveResponseBody"`       // Whether to save the response body in the monitor response
	SaveResponseHeaders  bool              `json:"saveResponseHeaders" bson:"saveResponseHeaders"` // Whether to save the response headers in the monitor response
	ExpectedStatusCodes  []int             `json:"expectedStatusCodes" bson:"expectedStatusCodes"`
	ExpectedBodyRegex    string            `json:"expectedBodyRegex" bson:"expectedBodyRegex"`
	ExpectedHeaders      map[string]string `json:"expectedHeaders" bson:"expectedHeaders"`
	ExpectedResponseTime *int              `json:"expectedResponseTime" bson:"expectedResponseTime"` // in milliseconds
}

type httpMonitor struct {
	BaseMonitor `bson:",inline"`
	Config      httpConfig `json:"config" bson:"config"`
}

func (m *httpMonitor) Run() IMonitorResponse {
	return m.Config.run()
}

func (m *httpMonitor) Validate() error {
	if err := m.validateBase(); err != nil {
		return fmt.Errorf("monitor validation failed: %w", err)
	}
	if err := m.Config.validate(); err != nil {
		return fmt.Errorf("HTTP monitor config validation failed: %w", err)
	}
	return nil
}

// httpClient is needed for mocking HTTP requests in tests
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func newHttpConfig(httpMethod, url string, headers map[string]string, body string, expectedStatusCodes []int, expectedBodyRegex string, expectedHeaders map[string]string, expectedResponseTime int) (*httpConfig, error) {
	monitor := &httpConfig{
		Method:               httpMethod,
		URL:                  url,
		Headers:              headers,
		Body:                 body,
		ExpectedStatusCodes:  expectedStatusCodes,
		ExpectedBodyRegex:    expectedBodyRegex,
		ExpectedHeaders:      expectedHeaders,
		ExpectedResponseTime: &expectedResponseTime,
	}

	if err := monitor.validate(); err != nil {
		return nil, err
	}

	return monitor, nil
}

func newHttpClient() httpClient {
	return &http.Client{
		Timeout: 10 * time.Second, // Default timeout for HTTP requests
	}
}

var httpClientOrMock = newHttpClient()

func (m *httpConfig) run() IMonitorResponse {
	monitorResponse := NewHttpMonitorResponse()

	httpResponse, elapsed, err := m.executeRequest(&httpClientOrMock)

	monitorResponse.Duration = elapsed.Milliseconds()
	if err != nil {
		// If the request failed altogether, there's no point in checking the httpResponse
		monitorResponse.setStatus(Error)
		monitorResponse.addErrorMsg(fmt.Sprintf("HTTP request failed: %s", err.Error()))
		logging.Uptime.Error().Err(err).Msg("HTTP monitor validation failed")
		return monitorResponse
	}

	monitorResponse.setRawHttpResponse(httpResponse, m.SaveResponseBody, m.SaveResponseHeaders)

	m.checkStatusCode(httpResponse, monitorResponse)
	m.checkResponseTime(elapsed, monitorResponse)
	m.checkResponseHeaders(httpResponse, monitorResponse)
	m.checkResponseBody(httpResponse, monitorResponse)

	return monitorResponse
}

// Encapsulates request creation and execution.
func (m *httpConfig) executeRequest(httpClient *httpClient) (*http.Response, time.Duration, error) {
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

func (m *httpConfig) checkStatusCode(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedStatusCodes == nil {
		return
	}

	if !util.SliceContains(m.ExpectedStatusCodes, response.StatusCode) {
		failureMsg := fmt.Sprintf("Unexpected status code: got %d, expected one of %d", response.StatusCode, m.ExpectedStatusCodes)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(statusCodeAspect)
	}
}

func (m *httpConfig) checkResponseTime(elapsed time.Duration, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedResponseTime == nil {
		return
	}
	if elapsed.Milliseconds() > int64(*m.ExpectedResponseTime) {
		failureMsg := fmt.Sprintf("Response time exceeded: got %dms, expected <= %dms", elapsed.Milliseconds(), *m.ExpectedResponseTime)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(responseTimeAspect)
	}
}

func (m *httpConfig) checkResponseHeaders(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if len(m.ExpectedHeaders) == 0 {
		return
	}

	for key, expectedValue := range m.ExpectedHeaders {
		actualValue := response.Header.Get(key)
		if actualValue != expectedValue {
			failureMsg := fmt.Sprintf("Header mismatch for %s: got %s, expected %s", key, actualValue, expectedValue)
			monitorResponse.addFailureMsg(failureMsg)
			monitorResponse.addFailedAspect(headersAspect)
		}
	}
}

func (m *httpConfig) checkResponseBody(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedBodyRegex == "" {
		return
	}

	responseBody, err := readResponseBody(response)
	if err != nil {
		monitorResponse.addErrorMsg("Error reading response body: " + err.Error())
		return
	}

	// Add (?s) flag to make dot match newlines
	patternWithFlag := "(?s)" + m.ExpectedBodyRegex

	regex, err := regexp.Compile(patternWithFlag)
	if err != nil {
		monitorResponse.addErrorMsg(fmt.Sprintf("Invalid regex for expected body: %s", patternWithFlag))
		return
	}

	matches := regex.Match([]byte(responseBody))
	if !matches {
		failureMsg := fmt.Sprintf("Response body does not match regex: %s", m.ExpectedBodyRegex)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(bodyAspect)
	}
}

// createRequest constructs an HTTP request based on the monitor's configuration
func (m *httpConfig) createRequest() (*http.Request, error) {
	parsedUrl, err := url.Parse(m.URL)
	if err != nil {
		logging.Uptime.Error().Err(err).Msg("Invalid URL in HTTP monitor")
		return nil, fmt.Errorf("invalid URL: %s", m.URL)
	}

	req := http.Request{
		Method: m.Method,
		URL:    parsedUrl,
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

// Validate checks if the HTTP monitor configuration is valid
// It ensures that required fields are set and that the URL is properly formatted.
func (m *httpConfig) validate() error {
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

	parsedUrl, err := url.Parse(m.URL)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return fmt.Errorf("invalid URL format: %s", m.URL)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return fmt.Errorf("URL scheme must be either http or https: %s", m.URL)
	}

	if m.ExpectedResponseTime != nil && *m.ExpectedResponseTime < 0 {
		return fmt.Errorf("expected response time cannot be negative")
	}

	return nil
}

// Helper function to read response body while preserving it
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

func (b *HttpMonitorResponse) setStatus(status monitorResponseStatus) {
	b.Status = status
}

func (b *HttpMonitorResponse) addFailedAspect(aspect httpCheckAspect) {
	if b.FailedAspects == nil {
		b.FailedAspects = make([]httpCheckAspect, 0)
	}
	b.FailedAspects = append(b.FailedAspects, aspect)
	if b.Status != Error {
		b.Status = failure
	}
}

func (m *httpMonitor) GetConfig() IMonitorConfig {
	return &m.Config
}

func (m *httpMonitor) SetConfig(config IMonitorConfig) {
	m.Config = *config.(*httpConfig)
}
