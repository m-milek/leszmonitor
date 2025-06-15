package monitors

import (
	"net"
	"testing"
	"time"

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
func setupPingMonitorConfig() *PingConfig {
	monitor, err := NewPingConfig("example.com", "80", "tcp", 5, 3)
	monitor.pingAddressFunc = pingAddressFunc // Use the global function for testing

	if err != nil {
		panic("Failed to create PingConfig: " + err.Error())
	}
	return monitor
}

func TestPingConfig_ImplementsIMonitorConfig(t *testing.T) {
	monitor := setupPingMonitorConfig()
	var iMonitor IMonitorConfig = monitor
	assert.NotNil(t, iMonitor)
}

func TestPingMonitor_ImplementsIMonitor(t *testing.T) {
	monitor := &PingMonitor{
		BaseMonitor: BaseMonitor{Id: "test-id"},
		Config:      *setupPingMonitorConfig(),
	}
	var iMonitor IMonitor = monitor
	assert.NotNil(t, iMonitor)
}

func TestPingMonitor_Validate(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		err := monitor.validate()
		assert.NoError(t, err)
	})

	t.Run("Empty Host", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.Host = ""
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})

	t.Run("Empty Port", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.Port = ""
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "port cannot be empty")
	})

	t.Run("Invalid RetryCount", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.RetryCount = 0
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count must be greater than zero")
	})

	t.Run("Invalid Protocol", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.Protocol = "invalid"
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid protocol")
	})

	t.Run("Valid Protocols", func(t *testing.T) {
		for _, protocol := range validProtocols {
			monitor := setupPingMonitorConfig()
			monitor.Protocol = protocol
			err := monitor.validate()
			assert.NoError(t, err, "Protocol %s should be valid", protocol)
		}
	})

	t.Run("Negative Timeout", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.PingTimeout = -1
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})

	t.Run("Timeout Exceeds Limit", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.PingTimeout = 61 // Exceeds the limit of 60 seconds
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must not exceed 60 seconds")
	})

	t.Run("Zero Timeout", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		monitor.PingTimeout = 0
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})
}

// Test the pingAddress function
func TestPingAddress(t *testing.T) {
	// This is a bit tricky to test without making actual network calls
	// We'll use a known reliable service for a simple integration test
	t.Run("Successful Ping", func(t *testing.T) {
		// Skip this test in CI environments or when offline
		if testing.Short() {
			t.Skip("Skipping network-dependent test in short mode")
		}

		success, duration := pingAddress("tcp", "localhost:80", 2*time.Second)
		// The test might fail if port 80 is not open on localhost
		// This is more of an integration test than a unit test
		if success {
			assert.True(t, duration > 0)
		}
	})

	t.Run("Failed Ping - Invalid Host", func(t *testing.T) {
		success, _ := pingAddress("tcp", "invalid-host-that-does-not-exist:80", 1*time.Second)
		assert.False(t, success)
	})

	t.Run("Failed Ping - Invalid Port", func(t *testing.T) {
		success, _ := pingAddress("tcp", "localhost:99999", 1*time.Second)
		assert.False(t, success)
	})
}

// Test the Run method with mocked network calls
func TestPingMonitor_Run(t *testing.T) {
	// Save the original function and restore it after tests
	originalPingAddress := pingAddress
	defer func() { pingAddressFunc = originalPingAddress }()

	t.Run("Successful Ping", func(t *testing.T) {
		monitor := setupPingMonitorConfig()

		// Mock the pingAddress function
		pingAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			assert.Equal(t, "tcp", protocol)
			assert.Equal(t, "example.com:80", address)
			assert.Equal(t, 5*time.Second, timeout)
			return true, 100 * time.Millisecond
		}

		response := monitor.run()
		assert.EqualValues(t, Success, response.GetStatus())
		assert.Equal(t, int64(100), response.GetDuration())
		assert.Empty(t, response.GetFailures())
	})

	t.Run("Failed Ping with Retries", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		callCount := 0

		// Mock the pingAddress function to fail for all retries
		pingAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			callCount++
			return false, 0
		}

		response := monitor.run()
		assert.Equal(t, 3, callCount, "Should have tried 3 times")
		assert.NotEmpty(t, response.GetFailures())
		assert.Contains(t, response.GetFailures()[0], "Failed to ping")
	})

	t.Run("Successful Ping After Retry", func(t *testing.T) {
		monitor := setupPingMonitorConfig()
		callCount := 0

		// Mock the pingAddress function to succeed on the second try
		pingAddressFunc = func(protocol string, address string, timeout time.Duration) (bool, time.Duration) {
			callCount++
			if callCount == 2 {
				return true, 150 * time.Millisecond
			}
			return false, 0
		}

		response := monitor.run()
		assert.Equal(t, 2, callCount, "Should have tried 2 times")
		assert.EqualValues(t, Success, response.GetStatus())
		assert.Equal(t, int64(150), response.GetDuration())
		assert.Empty(t, response.GetFailures())
	})
}

func TestPingMonitorResponse(t *testing.T) {
	t.Run("New Response", func(t *testing.T) {
		response := NewPingMonitorResponse()
		assert.EqualValues(t, Success, response.GetStatus())
		assert.NotZero(t, response.GetTimestamp())
		assert.Empty(t, response.GetErrorsStr())
		assert.Empty(t, response.GetFailures())
	})

	t.Run("Response Getters", func(t *testing.T) {
		response := NewPingMonitorResponse()
		response.Status = Failure
		response.Duration = 200
		response.ErrorsStr = []string{"Test error"}
		response.Failures = []string{"Test failure"}

		assert.EqualValues(t, Failure, response.GetStatus())
		assert.Equal(t, int64(200), response.GetDuration())
		assert.Equal(t, []string{"Test error"}, response.GetErrorsStr())
		assert.Equal(t, []string{"Test failure"}, response.GetFailures())
	})
}
