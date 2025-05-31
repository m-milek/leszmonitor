package monitors

import (
	"fmt"
	"github.com/m-milek/leszmonitor/util"
	"net"
	"time"
)

var (
	validProtocols = []string{
		"tcp",  // Transmission Control Protocol
		"udp",  // User Datagram Protocol
		"tcp4", // IPv4 over TCP
		"tcp6", // IPv6 over TCP
		"udp4", // IPv4 over UDP
		"udp6", // IPv6 over UDP
	}
	retryTimeout = 1 * time.Second // Default retry timeout
)

type PingMonitor struct {
	Base            baseMonitor `json:"base" bson:"base,inline"`
	Host            string      `json:"host" bson:"host"`               // Host to ping
	Port            string      `json:"port" bson:"port"`               // Port to ping
	Protocol        string      `json:"protocol" bson:"protocol"`       // Protocol to use (tcp, udp, etc.)
	Timeout         int         `json:"timeout" bson:"timeout"`         // Timeout in seconds for each ping
	RetryCount      int         `json:"retry_count" bson:"retry_count"` // Number of retries on failure
	pingAddressFunc func(protocol string, address string, timeout time.Duration) (bool, time.Duration)
}

func (m *PingMonitor) GetId() string {
	return m.Base.Id
}

func (m *PingMonitor) GetName() string {
	return m.Base.Name
}

func (m *PingMonitor) GetDescription() string {
	return m.Base.Description
}

func (m *PingMonitor) GetInterval() int {
	return m.Base.Interval
}

func (m *PingMonitor) GetType() MonitorType {
	return m.Base.Type
}

func (m *PingMonitor) setBase(base baseMonitor) {
	m.Base = base
}

func (m *PingMonitor) Run() (IMonitorResponse, error) {
	monitorResponse := NewPingMonitorResponse()

	// Handles IPv6 as well
	address := net.JoinHostPort(m.Host, m.Port)

	for i := 0; i < m.RetryCount; i++ {
		success, duration := pingAddressFunc(m.Protocol, address, time.Duration(m.Timeout)*time.Second)
		if success {
			monitorResponse.base.Duration = duration.Milliseconds()
			return monitorResponse, nil
		}
		if i < m.RetryCount-1 {
			time.Sleep(retryTimeout)
		}
	}

	// If we reach here, all retries failed
	monitorResponse.base.addFailureMsg(fmt.Sprintf("Failed to ping %s after %d retries", address, m.RetryCount))

	return monitorResponse, nil
}

func (m *PingMonitor) validateBase() error {
	return validateBaseMonitor(m)
}

func (m *PingMonitor) validate() error {
	baseErr := m.validateBase()
	if baseErr != nil {
		return baseErr
	}

	if m.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if m.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if m.RetryCount <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}

	if m.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}

	if m.Timeout > 60 {
		return fmt.Errorf("timeout must not exceed 60 seconds")
	}

	if !util.SliceContains(validProtocols, m.Protocol) {
		return fmt.Errorf("invalid protocol: %s, must be one of: %v", m.Protocol, validProtocols)
	}

	return nil
}

// pingAddressFunc is a function variable that can be replaced for testing purposes
var pingAddressFunc = pingAddress

// // pingAddress attempts to connect to the specified address using the given protocol
func pingAddress(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
	start := time.Now()
	conn, err := net.DialTimeout(protocol, address, timeout)
	duration := time.Since(start)

	if err != nil {
		return false, 0
	}

	defer conn.Close()
	return true, duration
}

func NewPingMonitor(base baseMonitor, host, port, protocol string, timeout, retryCount int) (*PingMonitor, error) {
	base.Type = Ping

	monitor := &PingMonitor{
		Base:            base,
		Host:            host,
		Port:            port,
		Protocol:        protocol,
		Timeout:         timeout,
		RetryCount:      retryCount,
		pingAddressFunc: pingAddress,
	}

	if err := monitor.validate(); err != nil {
		return nil, fmt.Errorf("failed to create PingMonitor: %w", err)
	}

	return monitor, nil
}

type PingMonitorResponse struct {
	base baseMonitorResponse `json:"base" bson:"base,inline"`
}

func NewPingMonitorResponse() *PingMonitorResponse {
	return &PingMonitorResponse{
		base: baseMonitorResponse{
			Status:    Success,
			Timestamp: util.GetUnixTimestamp(),
		},
	}
}

func (m *PingMonitorResponse) GetStatus() MonitorResponseStatus {
	return m.base.Status
}
func (m *PingMonitorResponse) GetDuration() int64 {
	return m.base.Duration
}
func (m *PingMonitorResponse) GetTimestamp() int64 {
	return m.base.Timestamp
}
func (m *PingMonitorResponse) GetErrors() []string {
	return m.base.Errors
}
func (m *PingMonitorResponse) GetFailures() []string {
	return m.base.Failures
}
