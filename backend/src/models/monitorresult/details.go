package monitorresult

type HttpResultDetails struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

type PingResultDetails struct {
	Retries   int `json:"retries"`
	LatencyMs int `json:"latencyMs"`
}

type IMonitorResultDetails interface {
	implDetails()
}

func (*HttpResultDetails) implDetails() {}

func (*PingResultDetails) implDetails() {}
