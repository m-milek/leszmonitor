package probe

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	icmpIDMask        = 0xffff
	icmpMaxPacketSize = 1500
)

type PingProbe struct {
	Host      string `json:"host" yaml:"host"`
	Count     int    `json:"count" yaml:"count"`
	TimeoutMs int    `json:"timeoutMs" yaml:"timeoutMs"`
}

type pingReply struct {
	rtt time.Duration
	err error
}

var (
	lookupIPFunc    = net.DefaultResolver.LookupIP
	pingAddressFunc = pingAddress
)

func (p *PingProbe) Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error) {
	details := &results.PingResultDetails{}
	result := results.NewMonitorResult(monitorID, kind.PingConfigType, kind.MonitorStatusUp, false, nil, details)

	ips, err := lookupIPFunc(ctx, "ip4", p.Host)
	if err != nil {
		result.AddFailure(
			results.FailureReasonPingNoReply,
			results.CauseFailureDetails{Cause: classifyNetError(err)},
			err,
		)
		return &result, nil
	}

	replies, err := pingAddressFunc(ctx, ips[0], p.Count, time.Duration(p.TimeoutMs)*time.Millisecond)
	if err != nil {
		return nil, err
	}

	var totalDuration time.Duration
	var lastErr error

	for _, reply := range replies {
		details.PacketsSent++
		if reply.err != nil {
			lastErr = reply.err
			continue
		}
		details.PacketsReceived++
		totalDuration += reply.rtt
	}

	if details.PacketsReceived < p.Count {
		result.AddFailure(
			results.FailureReasonPingNoReply,
			results.CauseFailureDetails{Cause: classifyNetError(lastErr)},
			lastErr,
		)
	}

	if details.PacketsReceived > 0 {
		avgDuration := totalDuration / time.Duration(details.PacketsReceived)
		details.LatencyMs = avgDuration.Milliseconds()
		result.SetDuration(new(avgDuration.Milliseconds()))
	}

	return &result, nil
}

func pingAddress(ctx context.Context, ip net.IP, count int, timeout time.Duration) ([]pingReply, error) {
	socket, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	stopClosingOnCancel := context.AfterFunc(ctx, func() { _ = socket.Close() })
	defer stopClosingOnCancel()

	dst := &net.UDPAddr{IP: ip}
	replies := make([]pingReply, 0, count)

	for seq := range count {
		rtt, pingErr := pingOnce(socket, dst, seq, timeout)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		replies = append(replies, pingReply{rtt: rtt, err: pingErr})
	}

	return replies, nil
}

func pingOnce(socket *icmp.PacketConn, dst *net.UDPAddr, seq int, timeout time.Duration) (time.Duration, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & icmpIDMask,
			Seq:  seq,
			Data: []byte("leszmonitor"),
		},
	}
	b, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	if _, err = socket.WriteTo(b, dst); err != nil {
		return 0, err
	}

	if err = socket.SetReadDeadline(start.Add(timeout)); err != nil {
		return 0, err
	}

	reply := make([]byte, icmpMaxPacketSize)
	for {
		n, _, readErr := socket.ReadFrom(reply)
		if readErr != nil {
			return 0, readErr
		}
		rm, parseErr := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
		if parseErr != nil {
			continue
		}
		if echo, ok := rm.Body.(*icmp.Echo); ok && rm.Type == ipv4.ICMPTypeEchoReply && echo.Seq == seq {
			return time.Since(start), nil
		}
	}
}

func (p *PingProbe) Validate() error {
	if p.Host == "" {
		return errors.New("host cannot be empty")
	}

	if p.TimeoutMs <= 0 {
		return fmt.Errorf("invalid ping timeout: %d, must be greater than zero", p.TimeoutMs)
	}

	if p.Count <= 0 {
		return fmt.Errorf("invalid ping count: %d, must be greater than zero", p.Count)
	}

	return nil
}
