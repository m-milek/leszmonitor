package monitors

import (
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/util"
	"net/http"
)

type httpCheckAspect string

const (
	statusCodeAspect   httpCheckAspect = "StatusCode"
	responseTimeAspect httpCheckAspect = "ResponseTime"
	bodyAspect         httpCheckAspect = "Body"
	headersAspect      httpCheckAspect = "Headers"
)

type HttpMonitorResponse struct {
	baseMonitorResponse `bson:",inline"`
	RawHttpResponse     RawHttpResponse   `json:"rawResponse" bson:"rawResponse"`
	FailedAspects       []httpCheckAspect `json:"failedAspects" bson:"failedAspects"` // Aspects that failed during the check
}

func NewHttpMonitorResponse() *HttpMonitorResponse {
	return &HttpMonitorResponse{
		baseMonitorResponse: baseMonitorResponse{
			Status:    Success,
			Timestamp: util.GetUnixTimestamp(),
			Failures:  []string{},
			ErrorsStr: []string{},
			Errors:    []error{},
		},
		FailedAspects: []httpCheckAspect{},
	}
}

type RawHttpResponse struct {
	StatusCode    int               `json:"statusCode" bson:"statusCode"`       // HTTP status code of the response
	Headers       map[string]string `json:"headers" bson:"headers"`             // Headers of the response
	Body          string            `json:"body" bson:"body"`                   // Body of the response
	ContentLength int64             `json:"contentLength" bson:"contentLength"` // Content length of the response body
	Proto         string            `json:"proto" bson:"proto"`                 // Protocol used for the response (e.g., HTTP/1.1)
	Cookies       []*http.Cookie    `json:"cookies" bson:"cookies"`             // Cookies set in the response
}

func NewRawHttpResponse(resp *http.Response, saveBody bool, saveHeaders bool) *RawHttpResponse {
	var body string
	if saveBody {
		readBody, err := readResponseBody(resp)

		if err != nil {
			logger.Uptime.Warn().Err(err).Msg("Failed to read HTTP response body")
			body = ""
		}
		body = readBody
	}

	var headers map[string]string
	if saveHeaders {
		headers = make(map[string]string)
		for key, value := range resp.Header {
			// Join multiple header values with commas
			headers[key] = value[0]
			if len(value) > 1 {
				headers[key] += ", " + value[1]
			}
		}
	}

	return &RawHttpResponse{
		StatusCode:    resp.StatusCode,
		Headers:       headers,
		Body:          body,
		ContentLength: resp.ContentLength,
		Proto:         resp.Proto,
		Cookies:       resp.Cookies(),
	}
}

func (b *HttpMonitorResponse) setRawHttpResponse(resp *http.Response, saveBody bool, saveHeaders bool) {
	b.RawHttpResponse = *NewRawHttpResponse(resp, saveBody, saveHeaders)
}
