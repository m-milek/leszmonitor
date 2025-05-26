package monitors

import "net/http"

type HttpMonitorResponse struct {
	Base            baseMonitorResponse `json:"base" bson:"base,inline"`
	RawHttpResponse *http.Response      `json:"raw_response" bson:"raw_response"`
	FailedAspects   []httpCheckAspect   `json:"failed_aspects" bson:"failed_aspects"` // Aspects that failed during the check
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

type httpCheckAspect string

const (
	statusCodeAspect   httpCheckAspect = "StatusCode"
	responseTimeAspect httpCheckAspect = "ResponseTime"
	bodyAspect         httpCheckAspect = "Body"
	headersAspect      httpCheckAspect = "Headers"
)
