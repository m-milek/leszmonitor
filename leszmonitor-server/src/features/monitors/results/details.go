package results

type HTTPResultDetails struct {
	StatusCode    int               `json:"statusCode"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	ContentLength int64             `json:"contentLength"`
	Proto         string            `json:"proto"`
}

type TCPResultDetails struct {
	LatencyMs int64 `json:"latencyMs"`
}

type IMonitorResultDetails any

type DNSResultDetails struct {
	ResolvedRecords []any `json:"resolvedRecords,omitempty"`
}
