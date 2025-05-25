package monitors

import (
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
	"net/url"
	"time"
)

type HttpMonitor struct {
	Base                 BaseMonitor       `json:"base" bson:"base,inline"`
	HttpMethod           string            `json:"http_method" bson:"http_method"`
	Url                  string            `json:"url" bson:"url"`
	Headers              map[string]string `json:"headers" bson:"headers"`
	Body                 string            `json:"body" bson:"body"`
	ExpectedStatusCode   int               `json:"expected_status_code" bson:"expected_status_code"`
	ExpectedBodyRegex    string            `json:"expected_body_regex" bson:"expected_body_regex"`
	ExpectedHeaders      map[string]string `json:"expected_headers" bson:"expected_headers"`
	ExpectedResponseTime int               `json:"expected_response_time" bson:"expected_response_time"` // in milliseconds
}

func (m *HttpMonitor) Run() error {
	client := &http.Client{
		Timeout: time.Duration(m.Base.Timeout) * time.Second,
	}

	request, err := m.createRequest()

	if err != nil {
		logger.Uptime.Error().Msg("Error creating HTTP request")
	}

	response, err := client.Do(request)

	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Error while sending HTTP request")
		return err
	}

	logger.Uptime.Trace().Any("response", response).Msg("HTTP response")

	return nil
}

func (m *HttpMonitor) createRequest() (*http.Request, error) {
	parsedUrl, err := url.Parse(m.Url)
	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Invalid URL in HTTP monitor")
		return nil, err
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
		req.Body = http.NoBody // Set the body if needed
	}

	return &req, nil
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

func (m *HttpMonitor) GetType() string {
	return MonitorTypeHttp
}
