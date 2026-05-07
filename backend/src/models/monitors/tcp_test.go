package monitors

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"

	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDialer is a mock implementation of the dialer interface used for testing.
type mockDialer struct {
	mock.Mock
}

func (m *mockDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	args := m.Called(network, address, timeout)
	return args.Get(0).(net.Conn), args.Error(1)
}

// Mock connection to simulate network behavior.
type mockConn struct {
	mock.Mock
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *mockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *mockConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *mockConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

// Setup function for tests
func setupTCPMonitorConfig() *TCPConfig {
	monitor, err := NewTCPConfig("example.com", 80, "tcp", 5000, 3)
	monitor.dialAddressFunc = dialAddressFunc // Use the global function for testing

	if err != nil {
		panic("Failed to create TCPConfig: " + err.Error())
	}
	return monitor
}

func TestTCPConfig_ImplementsIMonitorConfig(t *testing.T) {
	monitor := setupTCPMonitorConfig()
	var iMonitor IMonitorConfig = monitor
	assert.NotNil(t, iMonitor)
}

func TestTCPMonitor_ImplementsIMonitor(t *testing.T) {
	monitor := &TCPMonitor{
		BaseMonitor: BaseMonitor{Slug: "test-id"},
		Config:      *setupTCPMonitorConfig(),
	}
	var iMonitor IMonitor = monitor
	assert.NotNil(t, iMonitor)
}

func TestTCPMonitor_Validate(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		err := monitor.validate()
		assert.NoError(t, err)
	})

	t.Run("Empty Host", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		monitor.Host = ""
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})

	t.Run("Invalid RetryCount", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		monitor.RetryCount = 0
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count must be greater than zero")
	})

	t.Run("Invalid Protocol", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		monitor.Protocol = "invalid"
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid protocol")
	})

	t.Run("Valid Protocols", func(t *testing.T) {
		for _, protocol := range validProtocols {
			monitor := setupTCPMonitorConfig()
			monitor.Protocol = protocol
			err := monitor.validate()
			assert.NoError(t, err, "Protocol %s should be valid", protocol)
		}
	})

	t.Run("Negative Timeout", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		monitor.Timeout = -1
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})

	t.Run("Timeout Exceeds Limit", func(t *testing.T) {
		pConfig := setupTCPMonitorConfig()
		pConfig.Timeout = 60001 // Exceeds the limit of 60 seconds (60000ms)
		monitor := &TCPMonitor{
			BaseMonitor: BaseMonitor{Slug: "test-slug", Name: "Test Name", Interval: 60, Type: shared.TCPConfigType},
			Config:      *pConfig,
		}
		err := monitor.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must not exceed 60 seconds")
	})

	t.Run("Zero Timeout", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		monitor.Timeout = 0
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})
}

// Test the dialAddress function
func TestTCPAddress(t *testing.T) {
	// This is a bit tricky to test without making actual network calls
	// We'll use a known reliable service for a simple integration test
	t.Run("Successful shared.TCPConfigType", func(t *testing.T) {
		// Skip this test in CI environments or when offline
		if testing.Short() {
			t.Skip("Skipping network-dependent test in short mode")
		}

		success, duration := dialAddress("tcp", "localhost:80", 2*time.Second)
		// The test might fail if port 80 is not open on localhost
		// This is more of an integration test than a unit test
		if success {
			assert.True(t, duration > 0)
		}
	})

	t.Run("Failed shared.TCPConfigType - Invalid Host", func(t *testing.T) {
		success, _ := dialAddress("tcp", "invalid-host-that-does-not-exist:80", 1*time.Second)
		assert.False(t, success)
	})

	t.Run("Failed shared.TCPConfigType - Invalid Port", func(t *testing.T) {
		success, _ := dialAddress("tcp", "localhost:99999", 1*time.Second)
		assert.False(t, success)
	})
}

// Test the Run method with mocked network calls
func TestTCPMonitor_Run(t *testing.T) {
	// Save the original function and restore it after tests
	originalDialAddress := dialAddress
	defer func() { dialAddressFunc = originalDialAddress }()

	t.Run("Successful shared.TCPConfigType", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()

		// Mock the dialAddress function
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			assert.Equal(t, "tcp", protocol)
			assert.Equal(t, "example.com:80", address)
			assert.Equal(t, 5000*time.Millisecond, timeout)
			return true, 100 * time.Millisecond
		}

		response := monitor.run(uuid.Nil, shared.TCPConfigType)
		assert.True(t, response.GetIsSuccess())
		assert.Equal(t, int64(100), response.GetDurationMs())
		assert.Empty(t, response.GetErrorDetails().ErrorMessage)
	})

	t.Run("Failed shared.TCPConfigType with Retries", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		callCount := 0

		// Mock the dialAddress function to fail for all retries
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			callCount++
			return false, 0
		}

		response := monitor.run(uuid.Nil, shared.TCPConfigType)
		assert.Equal(t, 3, callCount, "Should have tried 3 times")
		assert.False(t, response.GetIsSuccess())
	})

	t.Run("Successful shared.TCPConfigType After Retry", func(t *testing.T) {
		monitor := setupTCPMonitorConfig()
		callCount := 0

		// Mock the dialAddress function to succeed on the second try
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			callCount++
			if callCount == 2 {
				return true, 150 * time.Millisecond
			}
			return false, 0
		}

		response := monitor.run(uuid.Nil, shared.TCPConfigType)
		assert.Equal(t, 2, callCount, "Should have tried 2 times")
		assert.True(t, response.GetIsSuccess())
		assert.Equal(t, int64(150), response.GetDurationMs())
	})
}
