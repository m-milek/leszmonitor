package probe

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

// Setup function for tests.
func setupTCPProbe() *TCPProbe {
	monitor, err := NewTCPProbe("example.com", 80, "tcp", 5000, 3)
	monitor.dialAddressFunc = dialAddressFunc // Use the global function for testing

	if err != nil {
		panic("Failed to create TCPConfig: " + err.Error())
	}
	return monitor
}

func TestTCPMonitor_Validate(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {
		probe := setupTCPProbe()
		err := probe.Validate()
		require.NoError(t, err)
	})

	t.Run("Empty Host", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.Host = ""
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})

	t.Run("Invalid RetryCount", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.RetryCount = 0
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "count must be greater than zero")
	})

	t.Run("Invalid Protocol", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.Protocol = "invalid"
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid protocol")
	})

	t.Run("Valid Protocols", func(t *testing.T) {
		for _, protocol := range validProtocols {
			probe := setupTCPProbe()
			probe.Protocol = protocol
			err := probe.Validate()
			assert.NoError(t, err, "Protocol %s should be valid", protocol)
		}
	})

	t.Run("Negative Timeout", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.Timeout = -1
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})

	t.Run("Timeout Exceeds Limit", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.Timeout = 60001 // Exceeds the limit of 60 seconds (60000ms)
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must not exceed 60 seconds")
	})

	t.Run("Zero Timeout", func(t *testing.T) {
		probe := setupTCPProbe()
		probe.Timeout = 0
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be greater than zero")
	})
}

// Test the dialAddress function.
func TestTCPAddress(t *testing.T) {
	// This is a bit tricky to test without making actual network calls
	// We'll use a known reliable service for a simple integration test
	t.Run("Successful TCPConfigType", func(t *testing.T) {
		// Skip this test in CI environments or when offline
		if testing.Short() {
			t.Skip("Skipping network-dependent test in short mode")
		}

		duration, err := dialAddress("tcp", "localhost:80", 2*time.Second)
		// The test might fail if port 80 is not open on localhost
		// This is more of an integration test than a unit test
		if err == nil {
			assert.Positive(t, duration)
		}
	})

	t.Run("Failed TCPConfigType - Invalid Host", func(t *testing.T) {
		_, err := dialAddress("tcp", "invalid-host-that-does-not-exist:80", 1*time.Second)
		assert.Error(t, err)
	})

	t.Run("Failed TCPConfigType - Invalid Port", func(t *testing.T) {
		_, err := dialAddress("tcp", "localhost:99999", 1*time.Second)
		assert.Error(t, err)
	})
}

// Test the Run method with mocked network calls.
func TestTCPMonitor_Run(t *testing.T) {
	// Save the original function and restore it after tests
	originalDialAddress := dialAddress
	defer func() { dialAddressFunc = originalDialAddress }()

	t.Run("Successful TCPConfigType", func(t *testing.T) {
		probe := setupTCPProbe()

		// Mock the dialAddress function
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (time.Duration, error) {
			assert.Equal(t, "tcp", protocol)
			assert.Equal(t, "example.com:80", address)
			assert.Equal(t, 5000*time.Millisecond, timeout)
			return 100 * time.Millisecond, nil
		}

		response, _ := probe.Run(context.Background(), uuid.Nil)
		assert.Equal(t, kind.MonitorStatusUp, response.GetStatus())
		assert.Equal(t, int64(100), response.GetDurationMs())
		assert.Empty(t, response.GetFailures())
	})

	t.Run("Failed TCPConfigType with Retries", func(t *testing.T) {
		probe := setupTCPProbe()
		callCount := 0

		// Mock the dialAddress function to fail for all retries
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (time.Duration, error) {
			callCount++
			return 0, syscall.ECONNREFUSED
		}

		response, _ := probe.Run(context.Background(), uuid.Nil)
		assert.Equal(t, 3, callCount, "Should have tried 3 times")
		assert.Equal(t, kind.MonitorStatusDown, response.GetStatus())
		assert.Equal(t, results.Failures{{
			Reason:  results.FailureReasonTCPConnectionFailed,
			Details: results.CauseFailureDetails{Cause: results.FailureCauseConnectionRefused},
			Error:   "connection refused",
		}}, response.GetFailures())
	})

	t.Run("Successful TCPConfigType After Retry", func(t *testing.T) {
		probe := setupTCPProbe()
		callCount := 0

		// Mock the dialAddress function to succeed on the second try
		dialAddressFunc = func(protocol string, address string, timeout time.Duration) (time.Duration, error) {
			callCount++
			if callCount == 2 {
				return 150 * time.Millisecond, nil
			}
			return 0, syscall.ECONNREFUSED
		}

		response, _ := probe.Run(context.Background(), uuid.Nil)
		assert.Equal(t, 2, callCount, "Should have tried 2 times")
		assert.Equal(t, kind.MonitorStatusUp, response.GetStatus())
		assert.Equal(t, int64(150), response.GetDurationMs())
	})
}
