package monitorresult

type HttpResultDetails struct {
	StatusCode    int               `json:"statusCode"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	ContentLength int64             `json:"contentLength"`
	Proto         string            `json:"proto"`
	FailedAspects []string          `json:"failedAspects,omitempty"`
}

type PingResultDetails struct {
	Tries     int64 `json:"tries"`
	LatencyMs int64 `json:"latencyMs"`
}

type IMonitorResultDetails interface {
	implDetails()
}

func (*HttpResultDetails) implDetails() {}

func (*PingResultDetails) implDetails() {}
