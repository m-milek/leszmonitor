package monitors

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	consts "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/util"
)

var (
	validProtocols = []string{
		"tcp",  // Transmission Control Protocol
		"tcp4", // IPv4 over TCP
		"tcp6", // IPv6 over TCP
	}
	retryTimeout = 1 * time.Second // Default retry timeout
)

type TCPProbe struct {
	Host            string `json:"host" bson:"host"`             // Host to call
	Port            int    `json:"port" bson:"port"`             // Port to call
	Protocol        string `json:"protocol" bson:"protocol"`     // Protocol to use (tcp, udp, etc.)
	Timeout         int    `json:"timeout" bson:"timeout"`       // Timeout in milliseconds for each connection attempt
	RetryCount      int    `json:"retryCount" bson:"retryCount"` // RetryCount is the number of retries until
	dialAddressFunc func(protocol string, address string, timeout time.Duration) (bool, time.Duration)
}

func (m *TCPProbe) Run(monitorID uuid.UUID) monitorresult.IMonitorResult {
	result := monitorresult.NewMonitorResult(monitorID, consts.TCPConfigType, true, false, 0, "", &monitorresult.TCPResultDetails{}, time.Now().Format(time.RFC3339))
	details := result.GetDetails().(*monitorresult.TCPResultDetails)

	portString := strconv.Itoa(m.Port)
	address := net.JoinHostPort(m.Host, portString)

	details.Tries++
	for i := 0; i < m.RetryCount; i++ {
		success, duration := dialAddressFunc(m.Protocol, address, time.Duration(m.Timeout)*time.Millisecond)
		if success {
			result.SetDuration(duration.Milliseconds())
			details.LatencyMs = duration.Milliseconds()
			return result
		}
		if i < m.RetryCount-1 {
			details.Tries++
			time.Sleep(retryTimeout)
		}
	}

	result.AddFailure(fmt.Sprintf("Failed to connect to %s after %d tries", address, m.RetryCount))

	return result
}

func (m *TCPProbe) Validate() error {
	if m.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if m.RetryCount <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}

	if m.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}

	// Timeout is in milliseconds; enforce a 60-second maximum
	if m.Timeout > 60000 {
		return fmt.Errorf("timeout must not exceed 60 seconds")
	}

	if !util.SliceContains(validProtocols, m.Protocol) {
		return fmt.Errorf("invalid protocol: %s, must be one of: %v", m.Protocol, validProtocols)
	}

	return nil
}

func NewTCPProbe(host string, port int, protocol string, timeout, retryCount int) (*TCPProbe, error) {
	probe := &TCPProbe{
		Host:            host,
		Port:            port,
		Protocol:        protocol,
		Timeout:         timeout,
		RetryCount:      retryCount,
		dialAddressFunc: dialAddress,
	}

	if err := probe.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create TCPConfig: %w", err)
	}

	return probe, nil
}

// dialAddressFunc is a function variable that can be replaced for testing purposes.
var dialAddressFunc = dialAddress

// // dialAddress attempts to connect to the specified address using the given protocol.
func dialAddress(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
	start := time.Now()
	conn, err := net.DialTimeout(protocol, address, timeout)
	duration := time.Since(start)

	if err != nil {
		return false, 0
	}

	defer conn.Close()
	return true, duration
}
