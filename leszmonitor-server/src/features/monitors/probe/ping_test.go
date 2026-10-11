package probe

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPingProbe() *PingProbe {
	return &PingProbe{Host: "example.com", Count: 3, TimeoutMs: 1000}
}

func mockPingNetwork(t *testing.T, replies []pingReply, pingErr error) *int {
	t.Helper()

	originalLookupIP := lookupIPFunc
	originalPingAddress := pingAddressFunc
	t.Cleanup(func() {
		lookupIPFunc = originalLookupIP
		pingAddressFunc = originalPingAddress
	})

	lookupIPFunc = func(_ context.Context, network, host string) ([]net.IP, error) {
		assert.Equal(t, "ip4", network)
		assert.Equal(t, "example.com", host)
		return []net.IP{net.IPv4(192, 0, 2, 1)}, nil
	}

	calls := 0
	pingAddressFunc = func(_ context.Context, ip net.IP, count int, timeout time.Duration) ([]pingReply, error) {
		calls++
		assert.True(t, ip.Equal(net.IPv4(192, 0, 2, 1)))
		assert.Equal(t, 3, count)
		assert.Equal(t, time.Second, timeout)
		return replies, pingErr
	}

	return &calls
}

func timeoutErr() error {
	return &net.OpError{Op: "read", Net: "udp4", Err: os.ErrDeadlineExceeded}
}

func TestPingProbe_Validate(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {
		require.NoError(t, setupPingProbe().Validate())
	})

	t.Run("Empty Host", func(t *testing.T) {
		probe := setupPingProbe()
		probe.Host = ""
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host cannot be empty")
	})

	t.Run("Zero Timeout", func(t *testing.T) {
		probe := setupPingProbe()
		probe.TimeoutMs = 0
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ping timeout")
	})

	t.Run("Zero Count", func(t *testing.T) {
		probe := setupPingProbe()
		probe.Count = 0
		err := probe.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ping count")
	})
}

func TestPingProbe_Run(t *testing.T) {
	t.Run("All Replies", func(t *testing.T) {
		mockPingNetwork(t, []pingReply{
			{rtt: 10 * time.Millisecond},
			{rtt: 20 * time.Millisecond},
			{rtt: 30 * time.Millisecond},
		}, nil)

		result, err := setupPingProbe().Run(context.Background(), uuid.Nil)
		require.NoError(t, err)

		assert.Equal(t, kind.MonitorStatusUp, result.GetStatus())
		assert.Empty(t, result.GetFailures())
		require.NotNil(t, result.GetDurationMs())
		assert.Equal(t, int64(20), *result.GetDurationMs())
		assert.Equal(
			t,
			&results.PingResultDetails{LatencyMs: 20, PacketsSent: 3, PacketsReceived: 3},
			result.GetDetails(),
		)
	})

	t.Run("Some Replies Lost", func(t *testing.T) {
		mockPingNetwork(t, []pingReply{
			{rtt: 10 * time.Millisecond},
			{err: timeoutErr()},
			{rtt: 30 * time.Millisecond},
		}, nil)

		result, err := setupPingProbe().Run(context.Background(), uuid.Nil)
		require.NoError(t, err)

		assert.Equal(t, kind.MonitorStatusDown, result.GetStatus())
		require.Len(t, result.GetFailures(), 1)
		assert.Equal(t, results.FailureReasonPingNoReply, result.GetFailures()[0].Reason)
		assert.Equal(
			t,
			results.CauseFailureDetails{Cause: results.FailureCauseTimeout},
			result.GetFailures()[0].Details,
		)
		assert.Equal(
			t,
			&results.PingResultDetails{LatencyMs: 20, PacketsSent: 3, PacketsReceived: 2},
			result.GetDetails(),
		)
	})

	t.Run("No Replies", func(t *testing.T) {
		mockPingNetwork(t, []pingReply{
			{err: timeoutErr()},
			{err: timeoutErr()},
			{err: timeoutErr()},
		}, nil)

		result, err := setupPingProbe().Run(context.Background(), uuid.Nil)
		require.NoError(t, err)

		assert.Equal(t, kind.MonitorStatusDown, result.GetStatus())
		assert.Nil(t, result.GetDurationMs())
		require.Len(t, result.GetFailures(), 1)
		assert.Equal(t, results.FailureReasonPingNoReply, result.GetFailures()[0].Reason)
		assert.Equal(t, &results.PingResultDetails{PacketsSent: 3}, result.GetDetails())
	})

	t.Run("DNS Failure", func(t *testing.T) {
		calls := mockPingNetwork(t, nil, nil)
		lookupIPFunc = func(context.Context, string, string) ([]net.IP, error) {
			return nil, &net.DNSError{Err: "no such host", Name: "example.com", IsNotFound: true}
		}

		result, err := setupPingProbe().Run(context.Background(), uuid.Nil)
		require.NoError(t, err)

		assert.Equal(t, kind.MonitorStatusDown, result.GetStatus())
		require.Len(t, result.GetFailures(), 1)
		assert.Equal(t, results.CauseFailureDetails{Cause: results.FailureCauseDNS}, result.GetFailures()[0].Details)
		assert.Equal(t, 0, *calls)
	})

	t.Run("Socket Error", func(t *testing.T) {
		mockPingNetwork(t, nil, os.ErrPermission)

		result, err := setupPingProbe().Run(context.Background(), uuid.Nil)
		require.ErrorIs(t, err, os.ErrPermission)
		assert.Nil(t, result)
	})
}
