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
	Hostname       string        `json:"hostname"`
	DNSServer      string        `json:"dnsServer,omitempty"`
	RecordType     DNSRecordType `json:"recordType"`
	ExpectedValues any           `json:"expectedValues"`
}

type DNSMXExpectedRecord struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}

type DNSSRVExpectedRecord struct {
	Target   string `json:"target"`
	Port     uint16 `json:"port"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
}

type DNSAExpectedValues []string
type DNSAAAAExpectedValues []string
type DNSTXTExpectedValues []string
type DNSMXExpectedValues []DNSMXExpectedRecord
type DNSNSExpectedValues []string
type DNSSRVExpectedValues []DNSSRVExpectedRecord

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
func checkExpected[E any, R any](
	result monitorresult.IMonitorResult,
	logger *zerolog.Logger,
	details *monitorresult.DNSResultDetails,
	expectedValues []E,
	resolvedRecords []R,
	recordToAny func(R) any,
	matches func(R, E) bool,
	notFoundMsg func(E) string,
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

	switch p.RecordType {
	case DNSRecordTypeA:
		expectedValues, ok := p.ExpectedValues.(DNSAExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for A record type", "Expected values type assertion failed for A record type")
		}
		ips, err := resolver.LookupIP(ctx, "ip4", p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup A records: %s", err.Error()), err, "A record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, ips,
			func(ip net.IP) any { return ip.String() },
			func(ip net.IP, expected string) bool { return ip.String() == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected A record with value %s not found", expected)
			},
		)

	case DNSRecordTypeAAAA:
		expectedValues, ok := p.ExpectedValues.(DNSAAAAExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for AAAA record type", "Expected values type assertion failed for AAAA record type")
		}
		ips, err := resolver.LookupIP(ctx, "ip6", p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup AAAA records: %s", err.Error()), err, "AAAA record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, ips,
			func(ip net.IP) any { return ip.String() },
			func(ip net.IP, expected string) bool { return ip.String() == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected AAAA record with value %s not found", expected)
			},
		)

	case DNSRecordTypeCNAME:
		expectedValues, ok := p.ExpectedValues.(DNSCNAMEExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for CNAME record type", "Expected values type assertion failed for CNAME record type")
		}
		cname, err := resolver.LookupCNAME(ctx, p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup CNAME record: %s", err.Error()), err, "CNAME record lookup failed")
		}
		details.ResolvedRecords = append(details.ResolvedRecords, cname)
		if cname != expectedValues.CNAME {
			msg := fmt.Sprintf("Expected CNAME record with value %s not found", expectedValues.CNAME)
			result.AddError(msg)
			logger.Trace().Msg(msg)
		}

	case DNSRecordTypeMX:
		expectedValues, ok := p.ExpectedValues.(DNSMXExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for MX record type", "Expected values type assertion failed for MX record type")
		}
		mxRecords, err := resolver.LookupMX(ctx, p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup MX records: %s", err.Error()), err, "MX record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, mxRecords,
			func(mx *net.MX) any { return mx },
			func(mx *net.MX, expected DNSMXExpectedRecord) bool {
				return mx.Host == expected.Host && mx.Pref == expected.Priority
			},
			func(expected DNSMXExpectedRecord) string {
				return fmt.Sprintf("Expected MX record with host %s and priority %d not found", expected.Host, expected.Priority)
			},
		)

	case DNSRecordTypeTXT:
		expectedValues, ok := p.ExpectedValues.(DNSTXTExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for TXT record type", "Expected values type assertion failed for TXT record type")
		}
		txtRecords, err := resolver.LookupTXT(ctx, p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup TXT records: %s", err.Error()), err, "TXT record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, txtRecords,
			func(txt string) any { return txt },
			func(txt string, expected string) bool { return txt == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected TXT record with value %s not found", expected)
			},
		)

	case DNSRecordTypeNS:
		expectedValues, ok := p.ExpectedValues.(DNSNSExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for NS record type", "Expected values type assertion failed for NS record type")
		}
		nsRecords, err := resolver.LookupNS(ctx, p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup NS records: %s", err.Error()), err, "NS record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, nsRecords,
			func(ns *net.NS) any { return ns },
			func(ns *net.NS, expected string) bool { return ns.Host == expected },
			func(expected string) string {
				return fmt.Sprintf("Expected NS record with value %s not found", expected)
			},
		)

	case DNSRecordTypeSRV:
		expectedValues, ok := p.ExpectedValues.(DNSSRVExpectedValues)
		if !ok {
			return earlyError(result, logger, "Invalid expected values for SRV record type", "Expected values type assertion failed for SRV record type")
		}
		service, proto, name, err := splitSRVHostname(p.Hostname)
		if err != nil {
			return earlyErrorWithErr(result, logger, err.Error(), err, "SRV record name parsing failed")
		}
		_, srvRecords, err := resolver.LookupSRV(ctx, service, proto, name)
		if err != nil {
			return earlyErrorWithErr(result, logger, fmt.Sprintf("Failed to lookup SRV records: %s", err.Error()), err, "SRV record lookup failed")
		}
		checkExpected(result, logger, details, expectedValues, srvRecords,
			func(srv *net.SRV) any { return srv },
			func(srv *net.SRV, expected DNSSRVExpectedRecord) bool { return srvMatchesExpected(srv, expected) },
			func(expected DNSSRVExpectedRecord) string {
				return fmt.Sprintf("Expected SRV record with target %s and port %d not found", expected.Target, expected.Port)
			},
		)

	default:
		result.AddError(fmt.Sprintf("Unsupported DNS record type: %s", rt))
		logger.Trace().Msgf("Unsupported DNS record type: %s", rt)
		result.SetDuration(0)
	}

	return result
}

func (p *DNSProbe) Validate() error {
	if p.Hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	if p.RecordType == "" {
		return fmt.Errorf("record type cannot be empty")
	}
	return nil
}

func makeResolver(DNSServer string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", DNSServer)
		},
	}
}

func srvMatchesExpected(record *net.SRV, expected DNSSRVExpectedRecord) bool {
	if record.Target != expected.Target {
		return false
	}
	if expected.Port != 0 && record.Port != expected.Port {
		return false
	}
	if expected.Priority != 0 && record.Priority != expected.Priority {
		return false
	}
	if expected.Weight != 0 && record.Weight != expected.Weight {
		return false
	}
	return true
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
