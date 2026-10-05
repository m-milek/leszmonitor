package kind

type ProbeType string

type HTTPResultDetailsType string

var (
	HTTPConfigType ProbeType = "http"
	TCPConfigType  ProbeType = "tcp"
	DNSConfigType  ProbeType = "dns"
	PushConfigType ProbeType = "push"
)
