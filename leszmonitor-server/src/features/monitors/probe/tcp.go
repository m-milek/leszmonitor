package probe

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/util"
)

var (
	validProtocols = []string{
		"tcp",  // Transmission Control Protocol
		"tcp4", // IPv4 over TCP
		"tcp6", // IPv6 over TCP
	}
)

type TCPProbe struct {
	Host            string `json:"host"`     // Host to call
	Port            int    `json:"port"`     // Port to call
	Protocol        string `json:"protocol"` // Protocol to use (tcp, udp, etc.)
	Timeout         int    `json:"timeout"`  // Timeout in milliseconds for the connection attempt
	dialAddressFunc func(protocol string, address string, timeout time.Duration) (time.Duration, error)
}

func NewTCPProbe(host string, port int, protocol string, timeout int) (*TCPProbe, error) {
	probe := &TCPProbe{
		Host:            host,
		Port:            port,
		Protocol:        protocol,
		Timeout:         timeout,
		dialAddressFunc: dialAddress,
	}

	if err := probe.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create TCPConfig: %w", err)
	}

	return probe, nil
}

func (m *TCPProbe) Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error) {
	result := results.NewMonitorResult(
		monitorID,
		kind.TCPConfigType,
		kind.MonitorStatusUp,
		false,
		0,
		&results.TCPResultDetails{},
	)
	details := result.GetDetails().(*results.TCPResultDetails)

	portString := strconv.Itoa(m.Port)
	address := net.JoinHostPort(m.Host, portString)

	duration, err := dialAddressFunc(m.Protocol, address, time.Duration(m.Timeout)*time.Millisecond)
	if err != nil {
		result.AddFailure(
			results.FailureReasonTCPConnectionFailed,
			results.CauseFailureDetails{Cause: classifyNetError(err)},
			err,
		)
		return &result, nil
	}

	result.SetDuration(duration.Milliseconds())
	details.LatencyMs = duration.Milliseconds()

	return &result, nil
}

func (m *TCPProbe) Validate() error {
	if m.Host == "" {
		return fmt.Errorf("host cannot be empty")
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

// dialAddressFunc is a function variable that can be replaced for testing purposes.
var dialAddressFunc = dialAddress

// // dialAddress attempts to connect to the specified address using the given protocol.
func dialAddress(protocol string, address string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout(protocol, address, timeout)
	duration := time.Since(start)

	if err != nil {
		return 0, err
	}

	defer conn.Close()
	return duration, nil
}
