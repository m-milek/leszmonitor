package monitors

import (
	"fmt"
	"github.com/m-milek/leszmonitor/logger"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type HttpMonitor struct {
	Base                 BaseMonitor       `json:"base" bson:"base,inline"`
	HttpMethod           string            `json:"http_method" bson:"http_method"`
	Url                  string            `json:"url" bson:"url"`
	Headers              map[string]string `json:"headers" bson:"headers"`
	Body                 string            `json:"body" bson:"body"`
	ExpectedStatusCode   *int              `json:"expected_status_code" bson:"expected_status_code"`
	ExpectedBodyRegex    string            `json:"expected_body_regex" bson:"expected_body_regex"`
	ExpectedHeaders      map[string]string `json:"expected_headers" bson:"expected_headers"`
	ExpectedResponseTime *int              `json:"expected_response_time" bson:"expected_response_time"` // in milliseconds
}

type HttpMonitorResponse struct {
	Base            BaseMonitorResponse `json:"base" bson:"base,inline"`
	RawHttpResponse *http.Response      `json:"raw_response" bson:"raw_response"`
	FailedAspects   []HttpCheckAspect   `json:"failed_aspects" bson:"failed_aspects"` // Aspects that failed during the check
}

func (b *HttpMonitorResponse) GetStatus() MonitorResponseStatus {
	return b.Base.Status
}
func (b *HttpMonitorResponse) GetDuration() int64 {
	return b.Base.Duration
}
func (b *HttpMonitorResponse) GetTimestamp() int64 {
	return b.Base.Timestamp
}
func (b *HttpMonitorResponse) GetErrors() []string {
	return b.Base.Errors
}
func (b *HttpMonitorResponse) GetFailures() []string {
	return b.Base.Failures
}

type HttpCheckAspect string

const (
	StatusCodeAspect   HttpCheckAspect = "StatusCode"
	ResponseTimeAspect HttpCheckAspect = "ResponseTime"
	BodyAspect         HttpCheckAspect = "Body"
	HeadersAspect      HttpCheckAspect = "Headers"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHttpMonitor(base BaseMonitor, httpMethod, url string, headers map[string]string, body string, expectedStatusCode int, expectedBodyRegex string, expectedHeaders map[string]string, expectedResponseTime int) (*HttpMonitor, error) {
	monitor := &HttpMonitor{
		Base:                 base,
		HttpMethod:           httpMethod,
		Url:                  url,
		Headers:              headers,
		Body:                 body,
		ExpectedStatusCode:   &expectedStatusCode,
		ExpectedBodyRegex:    expectedBodyRegex,
		ExpectedHeaders:      expectedHeaders,
		ExpectedResponseTime: &expectedResponseTime,
	}

	if err := monitor.validate(); err != nil {
		return nil, err
	}

	return monitor, nil
}

func (m *HttpMonitor) Run(httpClient HttpClient) (IMonitorResponse, error) {
	monitorResponse := newHttpMonitorResponse()

	response, elapsed, err := m.executeRequest(httpClient)

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
func (m *HttpMonitor) executeRequest(httpClient HttpClient) (*http.Response, time.Duration, error) {
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

// getMismatchedHeaders checks the response headers against the expected headers
func (m *HttpMonitor) getMismatchedHeaders(response *http.Response) map[string]string {
	mismatchedHeaders := map[string]string{}
	if len(m.ExpectedHeaders) == 0 {
		return mismatchedHeaders
	}

	for key, expectedValue := range m.ExpectedHeaders {
		actualValue := response.Header.Get(key)
		if actualValue != expectedValue {
			mismatchedHeaders[key] = actualValue
		}
	}

	return mismatchedHeaders
}

func (m *HttpMonitor) validateStatusCode(response *http.Response, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedStatusCode == nil {
		return
	}
	if response.StatusCode != *m.ExpectedStatusCode {
		failureMsg := fmt.Sprintf("Unexpected status code: got %d, expected %d", response.StatusCode, *m.ExpectedStatusCode)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(StatusCodeAspect)
	}
}

func (m *HttpMonitor) validateResponseTime(elapsed time.Duration, monitorResponse *HttpMonitorResponse) {
	if m.ExpectedResponseTime == nil {
		return
	}
	if elapsed.Milliseconds() > int64(*m.ExpectedResponseTime) {
		failureMsg := fmt.Sprintf("Response time exceeded: got %dms, expected <= %dms", elapsed.Milliseconds(), *m.ExpectedResponseTime)
		monitorResponse.addFailureMsg(failureMsg)
		monitorResponse.addFailedAspect(ResponseTimeAspect)
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
			monitorResponse.addFailedAspect(HeadersAspect)
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
		monitorResponse.addFailedAspect(BodyAspect)
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

func newHttpClient() HttpClient {
	return &http.Client{
		Timeout: 10 * time.Second, // Default timeout for HTTP requests
	}
}

// Helper function to create a new response object with default values
func newHttpMonitorResponse() *HttpMonitorResponse {
	return &HttpMonitorResponse{
		Base: BaseMonitorResponse{
			Status: Success,
			Errors: []string{},
		},
		FailedAspects: []HttpCheckAspect{},
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
		return fmt.Errorf("URL cannot be empty in HTTP monitor")
	}

	if m.HttpMethod == "" {
		return fmt.Errorf("HTTP method cannot be empty in HTTP monitor")
	}

	if m.ExpectedStatusCode != nil && (*m.ExpectedStatusCode < 100 || *m.ExpectedStatusCode > 599) {
		return fmt.Errorf("expected status code must be between 100 and 599 in HTTP monitor")
	}

	parsedUrl, err := url.Parse(m.Url)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return fmt.Errorf("invalid URL format in HTTP monitor: %s", m.Url)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return fmt.Errorf("URL scheme must be either http or https in HTTP monitor: %s", m.Url)
	}

	if m.ExpectedResponseTime != nil && *m.ExpectedResponseTime < 0 {
		return fmt.Errorf("expected response time cannot be negative in HTTP monitor")
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

func (m *HttpMonitor) GetTimeout() int {
	return m.Base.Timeout
}

func (m *HttpMonitor) GetType() MonitorType {
	return Http
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

func (b *HttpMonitorResponse) addFailedAspect(aspect HttpCheckAspect) {
	if b.FailedAspects == nil {
		b.FailedAspects = make([]HttpCheckAspect, 0)
	}
	b.FailedAspects = append(b.FailedAspects, aspect)
	if b.Base.Status != Error {
		b.Base.Status = Failure
	}
}
