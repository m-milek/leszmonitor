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
		"udp",  // User Datagram Protocol
		"tcp4", // IPv4 over TCP
		"tcp6", // IPv6 over TCP
		"udp4", // IPv4 over UDP
		"udp6", // IPv6 over UDP
	}
	retryTimeout = 1 * time.Second // Default retry timeout
)

type PingConfig struct {
	Host            string `json:"host" bson:"host"`             // Host to pingType
	Port            int    `json:"port" bson:"port"`             // Port to pingType
	Protocol        string `json:"protocol" bson:"protocol"`     // Protocol to use (tcp, udp, etc.)
	PingTimeout     int    `json:"timeout" bson:"timeout"`       // PingTimeout in milliseconds for each pingType
	RetryCount      int    `json:"retryCount" bson:"retryCount"` // RetryCount is the number of retries until
	pingAddressFunc func(protocol string, address string, timeout time.Duration) (bool, time.Duration)
}

type PingMonitor struct {
	BaseMonitor `bson:",inline"` // Embed BaseMonitor for common fields
	Config      PingConfig       `json:"config" bson:"config"`
}

func (m *PingMonitor) GetConfig() IMonitorConfig {
	return &m.Config
}

func (m *PingMonitor) SetConfig(config IMonitorConfig) {
	m.Config = *config.(*PingConfig)
}

func (m *PingMonitor) Run() monitorresult.IMonitorResult {
	return m.Config.run(m.ID, m.Type)
}

func (m *PingMonitor) Validate() error {
	if err := m.validateBase(); err != nil {
		return fmt.Errorf("monitor validation failed: %w", err)
	}
	if err := m.Config.validate(); err != nil {
		return fmt.Errorf("pingType monitor config validation failed: %w", err)
	}
	return nil
}

func NewPingConfig(host string, port int, protocol string, timeout, retryCount int) (*PingConfig, error) {
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

func (m *PingConfig) run(id uuid.UUID, monitorType consts.MonitorConfigType) monitorresult.IMonitorResult {
	result := monitorresult.NewMonitorResult(id, monitorType, true, false, 0, "", &monitorresult.PingResultDetails{}, time.Now().Format(time.RFC3339))
	details := result.GetDetails().(*monitorresult.PingResultDetails)

	portString := strconv.Itoa(m.Port)
	address := net.JoinHostPort(m.Host, portString)

	details.Tries++
	for i := 0; i < m.RetryCount; i++ {
		success, duration := pingAddressFunc(m.Protocol, address, time.Duration(m.PingTimeout)*time.Millisecond)
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

	result.AddFailure(fmt.Sprintf("Failed to ping %s after %d tries", address, m.RetryCount))

	return result
}

func (m *PingConfig) validate() error {
	if m.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if m.RetryCount <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}

	if m.PingTimeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}

	// PingTimeout is in milliseconds; enforce a 60-second maximum
	if m.PingTimeout > 60000 {
		return fmt.Errorf("timeout must not exceed 60 seconds")
	}

	if !util.SliceContains(validProtocols, m.Protocol) {
		return fmt.Errorf("invalid protocol: %s, must be one of: %v", m.Protocol, validProtocols)
	}

	return nil
}

// pingAddressFunc is a function variable that can be replaced for testing purposes.
var pingAddressFunc = pingAddress

// // pingAddress attempts to connect to the specified address using the given protocol.
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
