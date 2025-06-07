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

type PingConfig struct {
	Host            string `json:"host" bson:"host"`               // Host to ping
	Port            string `json:"port" bson:"port"`               // Port to ping
	Protocol        string `json:"protocol" bson:"protocol"`       // Protocol to use (tcp, udp, etc.)
	PingTimeout     int    `json:"timeout" bson:"timeout"`         // PingTimeout in seconds for each ping
	RetryCount      int    `json:"retry_count" bson:"retry_count"` // Number of retries on failure
	pingAddressFunc func(protocol string, address string, timeout time.Duration) (bool, time.Duration)
}

type PingMonitor struct {
	BaseMonitor `bson:",inline"` // Embed BaseMonitor for common fields
	Config      PingConfig       `json:"config" bson:"config"`
}

func (m PingMonitor) Run() IMonitorResponse {
	return m.Config.run()
}

func (m *PingConfig) run() IMonitorResponse {
	monitorResponse := NewPingMonitorResponse()

	// Handles IPv6 as well
	address := net.JoinHostPort(m.Host, m.Port)

	for i := 0; i < m.RetryCount; i++ {
		success, duration := pingAddressFunc(m.Protocol, address, time.Duration(m.PingTimeout)*time.Second)
		if success {
			monitorResponse.Duration = duration.Milliseconds()
			return monitorResponse
		}
		if i < m.RetryCount-1 {
			time.Sleep(retryTimeout)
		}
	}

	// If we reach here, all retries failed
	monitorResponse.addFailureMsg(fmt.Sprintf("Failed to ping %s after %d retries", address, m.RetryCount))

	return monitorResponse
}

func (m *PingConfig) validate() error {
	if m.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if m.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if m.RetryCount <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}

	if m.PingTimeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}

	if m.PingTimeout > 60 {
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

func NewPingConfig(host, port, protocol string, timeout, retryCount int) (*PingConfig, error) {
	monitor := &PingConfig{
		Host:            host,
		Port:            port,
		Protocol:        protocol,
		PingTimeout:     timeout,
		RetryCount:      retryCount,
		pingAddressFunc: pingAddress,
	}

	if err := monitor.validate(); err != nil {
		return nil, fmt.Errorf("failed to create PingConfig: %w", err)
	}

	return monitor, nil
}

type PingMonitorResponse struct {
	baseMonitorResponse `bson:",inline"`
}

func NewPingMonitorResponse() *PingMonitorResponse {
	return &PingMonitorResponse{
		baseMonitorResponse: baseMonitorResponse{
			Status:    Success,
			Timestamp: util.GetUnixTimestamp(),
		},
	}
}

func (m *PingMonitorResponse) GetStatus() MonitorResponseStatus {
	return m.Status
}
func (m *PingMonitorResponse) GetDuration() int64 {
	return m.Duration
}
func (m *PingMonitorResponse) GetTimestamp() int64 {
	return m.Timestamp
}
func (m *PingMonitorResponse) GetErrors() []string {
	return m.Errors
}
func (m *PingMonitorResponse) GetFailures() []string {
	return m.Failures
}
func (m *PingMonitor) GenerateId() {
	if m.Id == "" {
		m.Id = generateMonitorId()
	}
}
