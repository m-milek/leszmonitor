package probe

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"syscall"

	"github.com/m-milek/leszmonitor/features/monitors/results"
)

func classifyNetError(err error) results.FailureCause {
	var dnsErr *net.DNSError
	var netErr net.Error
	var certErr *tls.CertificateVerificationError
	var recordErr tls.RecordHeaderError
	var alertErr tls.AlertError

	switch {
	case errors.As(err, &dnsErr):
		return results.FailureCauseDNS
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &netErr) && netErr.Timeout():
		return results.FailureCauseTimeout
	case errors.Is(err, syscall.ECONNREFUSED):
		return results.FailureCauseConnectionRefused
	case errors.Is(err, syscall.EHOSTUNREACH), errors.Is(err, syscall.ENETUNREACH):
		return results.FailureCauseUnreachable
	case errors.As(err, &certErr), errors.As(err, &recordErr), errors.As(err, &alertErr):
		return results.FailureCauseTLS
	default:
		return results.FailureCauseOther
	}
}

func classifyDNSError(err error) results.FailureCause {
	var dnsErr *net.DNSError
	switch {
	case errors.As(err, &dnsErr) && dnsErr.IsNotFound:
		return results.FailureCauseNotFound
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &dnsErr) && dnsErr.IsTimeout:
		return results.FailureCauseTimeout
	case errors.As(err, &dnsErr) && dnsErr.IsTemporary:
		return results.FailureCauseTemporary
	default:
		return results.FailureCauseOther
	}
}
