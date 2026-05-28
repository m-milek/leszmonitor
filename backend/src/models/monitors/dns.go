package monitors

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/rs/zerolog"
)

type DNSRecordType string

const (
	DNSRecordTypeA     DNSRecordType = "A"
	DNSRecordTypeAAAA  DNSRecordType = "AAAA"
	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeMX    DNSRecordType = "MX"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
	DNSRecordTypeNS    DNSRecordType = "NS"
	DNSRecordTypeSRV   DNSRecordType = "SRV"
)

type DNSProbe struct {
	Hostname             string        `json:"hostname"`
	DNSServer            string        `json:"dnsServer,omitempty"`
	RecordType           DNSRecordType `json:"recordType"`
	ExpectedRecordValues []string      `json:"expectedRecordValues"`
}

type DNSCNAMEExpectedValues struct {
	CNAME string `json:"cname"`
}

// earlyError sets duration to 0, adds the error message, logs it, and returns the result.
func earlyError(result monitorresult.IMonitorResult, logger *zerolog.Logger, userMsg, logMsg string) monitorresult.IMonitorResult {
	result.AddError(userMsg)
	logger.Trace().Msg(logMsg)
	result.SetDuration(0)
	return result
}

// earlyErrorWithErr is like earlyError but includes an error in the log.
func earlyErrorWithErr(result monitorresult.IMonitorResult, logger *zerolog.Logger, userMsg string, err error, logMsg string) monitorresult.IMonitorResult {
	result.AddError(userMsg)
	logger.Trace().Err(err).Msg(logMsg)
	result.SetDuration(0)
	return result
}

// checkExpected iterates over expected values and checks if each one is found in the resolved set.
// toStr converts a resolved record to a string for comparison.
// desc is used in error messages (e.g. "A", "AAAA", "TXT").
func checkExpected[R any](
	result monitorresult.IMonitorResult,
	logger *zerolog.Logger,
	details *monitorresult.DNSResultDetails,
	expectedValues []string,
	resolvedRecords []R,
	recordToAny func(R) any,
	matches func(R, string) bool,
	notFoundMsg func(string) string,
) {
	for _, r := range resolvedRecords {
		details.ResolvedRecords = append(details.ResolvedRecords, recordToAny(r))
	}

	for _, expected := range expectedValues {
		found := false
		for _, r := range resolvedRecords {
			if matches(r, expected) {
				found = true
				break
			}
		}
		if !found {
			msg := notFoundMsg(expected)
			result.AddFailure(msg)
			logger.Trace().Msg(msg)
		}
	}
}

func (p *DNSProbe) Run(ctx context.Context, monitorID uuid.UUID) monitorresult.IMonitorResult {
	logger := log.FromContext(ctx)
	result := monitorresult.NewMonitorResult(monitorID, consts.DNSConfigType, true, false, 0, "", &monitorresult.DNSResultDetails{})
	details := result.GetDetails().(*monitorresult.DNSResultDetails)

	resolver := makeResolver(p.DNSServer)
	rt := string(p.RecordType)

	startTime := time.Now()
	var endTime time.Time

	switch p.RecordType {
	case DNSRecordTypeA:
		ips, err := resolver.LookupIP(ctx, "ip4", p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup A records: %s", err.Error()), err, "A record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, ips,
			func(ip net.IP) any { return ip.String() },
			func(ip net.IP, expected string) bool { return ip.String() == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected A record with value %s not found", expected)
			},
		)

	case DNSRecordTypeAAAA:
		ips, err := resolver.LookupIP(ctx, "ip6", p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup AAAA records: %s", err.Error()), err, "AAAA record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, ips,
			func(ip net.IP) any { return ip.String() },
			func(ip net.IP, expected string) bool { return ip.String() == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected AAAA record with value %s not found", expected)
			},
		)

	case DNSRecordTypeCNAME:
		cname, err := resolver.LookupCNAME(ctx, p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup CNAME record: %s", err.Error()), err, "CNAME record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, []string{cname},
			func(cname string) any { return cname },
			func(cname string, expected string) bool { return cname == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected CNAME record with value %s not found", expected)
			},
		)

	case DNSRecordTypeMX:
		mxRecords, err := resolver.LookupMX(ctx, p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup MX records: %s", err.Error()), err, "MX record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, mxRecords,
			func(mx *net.MX) any { return fmt.Sprintf("%s:%d", mx.Host, mx.Pref) },
			func(mx *net.MX, expected string) bool { return fmt.Sprintf("%s:%d", mx.Host, mx.Pref) == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected MX record with value %s not found", expected)
			},
		)

	case DNSRecordTypeTXT:
		txtRecords, err := resolver.LookupTXT(ctx, p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup TXT records: %s", err.Error()), err, "TXT record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, txtRecords,
			func(txt string) any { return txt },
			func(txt string, expected string) bool { return txt == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected TXT record with value %s not found", expected)
			},
		)

	case DNSRecordTypeNS:
		nsRecords, err := resolver.LookupNS(ctx, p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup NS records: %s", err.Error()), err, "NS record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, nsRecords,
			func(ns *net.NS) any { return ns.Host },
			func(ns *net.NS, expected string) bool { return ns.Host == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected NS record with value %s not found", expected)
			},
		)

	case DNSRecordTypeSRV:
		service, proto, name, err := splitSRVHostname(p.Hostname)
		endTime = time.Now()
		if err != nil {
			return earlyErrorWithErr(result, logger, err.Error(), err, "SRV record name parsing failed")
		}
		_, srvRecords, err := resolver.LookupSRV(ctx, service, proto, name)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup SRV records: %s", err.Error()), err, "SRV record lookup failed")
		}
		checkExpected(result, logger, details, p.ExpectedRecordValues, srvRecords,
			func(srv *net.SRV) any { return fmt.Sprintf("%s:%d", srv.Target, srv.Port) },
			func(srv *net.SRV, expected string) bool {
				return fmt.Sprintf("%s:%d", srv.Target, srv.Port) == expected
			},
			func(expected string) string {
				return fmt.Sprintf("Expected SRV record with value %s not found", expected)
			},
		)

	default:
		result.AddError(fmt.Sprintf("Unsupported DNS record type: %s", rt))
		logger.Error().Msgf("Unsupported DNS record type: %s", rt)
		result.SetDuration(0)
	}

	result.SetDuration(endTime.Sub(startTime).Milliseconds())

	return result
}

func (p *DNSProbe) Validate() error {
	if p.Hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	if p.RecordType == "" {
		return fmt.Errorf("record type cannot be empty")
	}
	if p.ExpectedRecordValues == nil {
		return fmt.Errorf("expected values must be defined")
	}

	return nil
}

func makeResolver(dnsServer string) *net.Resolver {
	if dnsServer == "" {
		return net.DefaultResolver
	}

	// Handle case where port is already included
	host, port, err := net.SplitHostPort(dnsServer)
	if err != nil {
		// No port present, default to 53
		host = dnsServer
		port = "53"
	}
	addr := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", addr)
		},
	}
}

func splitSRVHostname(hostname string) (string, string, string, error) {
	parts := strings.SplitN(hostname, ".", 3)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid SRV hostname: %s", hostname)
	}

	service := strings.TrimPrefix(parts[0], "_")
	proto := strings.TrimPrefix(parts[1], "_")
	name := parts[2]
	if service == "" || proto == "" || name == "" {
		return "", "", "", fmt.Errorf("invalid SRV hostname: %s", hostname)
	}

	return service, proto, name, nil
}
