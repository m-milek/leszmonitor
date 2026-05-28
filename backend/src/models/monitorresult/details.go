package monitorresult

type HttpResultDetails struct {
	StatusCode    int               `json:"statusCode"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	ContentLength int64             `json:"contentLength"`
	Proto         string            `json:"proto"`
}

type TCPResultDetails struct {
	Tries     int64 `json:"tries"`
	LatencyMs int64 `json:"latencyMs"`
}

type IMonitorResultDetails interface {
	implDetails()
}

type DNSResultDetails struct {
	ResolvedRecords []any `json:"resolvedRecords,omitempty"`
}

func (*HttpResultDetails) implDetails() {}

func (*TCPResultDetails) implDetails() {}

func (*DNSResultDetails) implDetails() {}
