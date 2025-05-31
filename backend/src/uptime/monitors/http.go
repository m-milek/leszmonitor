package monitors

import (
	"fmt"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type HttpMonitor struct {
	Base                 baseMonitor       `json:"base,inline" bson:"base,inline"`
	HttpMethod           string            `json:"http_method" bson:"http_method"`
	Url                  string            `json:"url" bson:"url"`
	Headers              map[string]string `json:"headers" bson:"headers"`
	Body                 string            `json:"body" bson:"body"`
	ExpectedStatusCodes  []int             `json:"expected_status_codes" bson:"expected_status_codes"`
	ExpectedBodyRegex    string            `json:"expected_body_regex" bson:"expected_body_regex"`
	ExpectedHeaders      map[string]string `json:"expected_headers" bson:"expected_headers"`
	ExpectedResponseTime *int              `json:"expected_response_time" bson:"expected_response_time"` // in milliseconds
}

// httpClient is needed for mocking HTTP requests in tests
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHttpMonitor(base baseMonitor, httpMethod, url string, headers map[string]string, body string, expectedStatusCodes []int, expectedBodyRegex string, expectedHeaders map[string]string, expectedResponseTime int) (*HttpMonitor, error) {
	monitor := &HttpMonitor{
		Base:                 base,
		HttpMethod:           httpMethod,
		Url:                  url,
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

	base.Type = Http

	return monitor, nil
}

var httpClientInstance = newHttpClient()

func (m *HttpMonitor) Run() (IMonitorResponse, error) {
	monitorResponse := newHttpMonitorResponse()

	response, elapsed, err := m.executeRequest(httpClientInstance)

	monitorResponse.Base.Duration = elapsed.Milliseconds()
	if err != nil {
		// If the request failed altogether, there's no point in checking the response
		monitorResponse.setStatus(Error)
		monitorResponse.addErrorMsg(fmt.Sprintf("HTTP request failed: %s", err.Error()))
		logger.Uptime.Error().Err(err).Msg("HTTP monitor run failed")
		return monitorResponse, err
	}

	monitorResponse.RawHttpResponse = response

	m.validateStatusCode(response, monitorResponse)
	m.validateResponseTime(elapsed, monitorResponse)
	m.validateResponseHeaders(response, monitorResponse)
	m.validateResponseBody(response, monitorResponse)

	if monitorResponse.Base.Status == Error {
		errMsg := fmt.Sprintf("HTTP monitor run %s finished with errors: %v", m.Base.Name, monitorResponse.Base.Errors)
		logger.Uptime.Error().Any("monitor_response", monitorResponse).Msg(errMsg)
		return monitorResponse, fmt.Errorf("HTTP monitor run %s finished with errors: %v", m.Base.Name, monitorResponse.Base.Errors)
	}

	logger.Uptime.Debug().Any("monitor_response", monitorResponse).Msg("HTTP monitor run completed successfully")
	return monitorResponse, nil
}

// Encapsulates request creation and execution
func (m *HttpMonitor) executeRequest(httpClient httpClient) (*http.Response, time.Duration, error) {
	request, err := m.createRequest()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	start := time.Now()
	response, err := httpClient.Do(request)
	elapsed := time.Since(start)

	if err != nil {
		return nil, elapsed, err
	}

	return response, elapsed, nil
}

func (m *HttpMonitor) validateStatusCode(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedStatusCodes == nil {
		return
	}

	if !util.SliceContains(m.ExpectedStatusCodes, response.StatusCode) {
		failureMsg := fmt.Sprintf("Unexpected status code: got %d, expected one of %d", response.StatusCode, m.ExpectedStatusCodes)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(statusCodeAspect)
	}
}

func (m *HttpMonitor) validateResponseTime(elapsed time.Duration, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedResponseTime == nil {
		return
	}
	if elapsed.Milliseconds() > int64(*m.ExpectedResponseTime) {
		failureMsg := fmt.Sprintf("Response time exceeded: got %dms, expected <= %dms", elapsed.Milliseconds(), *m.ExpectedResponseTime)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(responseTimeAspect)
	}
}

func (m *HttpMonitor) validateResponseHeaders(response *http.Response, monitorResponse *HttpMonitorResponse) {
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

func (m *HttpMonitor) validateResponseBody(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedBodyRegex == "" {
		return
	}

	responseBody, err := readResponseBody(response)
	if err != nil {
		monitorResponse.addErrorMsg("Error reading response body: " + err.Error())
		return
	}

	matched, err := regexp.MatchString(m.ExpectedBodyRegex, responseBody)
	if err != nil {
		errMsg := fmt.Sprintf("Error matching response body regex: %s", err.Error())
		monitorResponse.addErrorMsg(errMsg)
		return
	}

	if !matched {
		failureMsg := fmt.Sprintf("Response body does not match regex: %s", m.ExpectedBodyRegex)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(bodyAspect)
	}
}

// createRequest constructs an HTTP request based on the monitor's configuration
func (m *HttpMonitor) createRequest() (*http.Request, error) {
	parsedUrl, err := url.Parse(m.Url)
	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Invalid URL in HTTP monitor")
		return nil, fmt.Errorf("invalid URL: %s", m.Url)
	}

	req := http.Request{
		Method: m.HttpMethod,
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

func newHttpClient() httpClient {
	return &http.Client{
		Timeout: 10 * time.Second, // Default timeout for HTTP requests
	}
}

// Helper function to create a new response object with default values
func newHttpMonitorResponse() *HttpMonitorResponse {
	return &HttpMonitorResponse{
		Base: baseMonitorResponse{
			Status: Success,
			Errors: []string{},
		},
		FailedAspects: []httpCheckAspect{},
	}
}

// validate checks if the HTTP monitor configuration is valid
// It ensures that required fields are set and that the URL is properly formatted.
func (m *HttpMonitor) validate() error {
	m.Base.Type = Http

	baseValidationErr := m.Base.validate()
	if baseValidationErr != nil {
		return baseValidationErr
	}

	if m.Url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if m.HttpMethod == "" {
		return fmt.Errorf("HTTP method cannot be empty")
	}

	if !util.SliceContains([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}, m.HttpMethod) {
		return fmt.Errorf("invalid HTTP method: %s", m.HttpMethod)
	}

	if m.ExpectedStatusCodes != nil && len(m.ExpectedStatusCodes) == 0 {
		return fmt.Errorf("expected status codes cannot be empty")
	}

	if m.ExpectedStatusCodes != nil && len(m.ExpectedStatusCodes) > 0 {
		minValue, maxValue := util.SliceMinMax(m.ExpectedStatusCodes)
		if minValue < 100 || maxValue > 599 {
			return fmt.Errorf("expected status codes must be between 100 and 599")
		}
	}

	parsedUrl, err := url.Parse(m.Url)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return fmt.Errorf("invalid URL format: %s", m.Url)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return fmt.Errorf("URL scheme must be either http or https: %s", m.Url)
	}

	if m.ExpectedResponseTime != nil && *m.ExpectedResponseTime < 0 {
		return fmt.Errorf("expected response time cannot be negative")
	}

	return nil
}

func readResponseBody(response *http.Response) (string, error) {
	if response.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Error reading response body")
		return "", err
	}

	if err := response.Body.Close(); err != nil {
		logger.Uptime.Error().Err(err).Msg("Error closing response body")
		return "", err
	}

	return string(bodyBytes), nil
}

func (m *HttpMonitor) GetName() string {
	return m.Base.Name
}

func (m *HttpMonitor) GetDescription() string {
	return m.Base.Description
}

func (m *HttpMonitor) GetInterval() int {
	return m.Base.Interval
}

func (m *HttpMonitor) GetType() MonitorType {
	return Http
}

func (m *HttpMonitor) setBase(base baseMonitor) {
	m.Base = base
}

func (b *HttpMonitorResponse) setStatus(status MonitorResponseStatus) {
	b.Base.Status = status
}

func (b *HttpMonitorResponse) addErrorMsg(err string) {
	b.Base.addErrorMsg(err)
}

func (b *HttpMonitorResponse) addFailureMsg(err string) {
	b.Base.addFailureMsg(err)
}

func (b *HttpMonitorResponse) addFailedAspect(aspect httpCheckAspect) {
	if b.FailedAspects == nil {
		b.FailedAspects = make([]httpCheckAspect, 0)
	}
	b.FailedAspects = append(b.FailedAspects, aspect)
	if b.Base.Status != Error {
		b.Base.Status = Failure
	}
}
