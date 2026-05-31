package monitors

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

const testDNSServer = "8.8.8.8:53"

func dnsProbe(hostname string, recordType DNSRecordType, expected ...string) DNSProbe {
	return DNSProbe{
		Hostname:             hostname,
		DNSServer:            testDNSServer,
		RecordType:           recordType,
		ExpectedRecordValues: expected,
	}
}

func assertDNSProbeRun(t *testing.T, probe DNSProbe, wantSuccess bool, timeout time.Duration) {
	t.Helper()

	ctx := t.Context()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result := probe.Run(ctx, uuid.New())

	if gotSuccess := result.GetIsSuccess(); gotSuccess != wantSuccess {
		if wantSuccess {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}
		t.Fatalf("Expected failure, got success")
	}

	if !wantSuccess {
		return
	}

	details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
	if !ok {
		t.Fatalf("Expected DNS result details")
	}

	if len(details.ResolvedRecords) == 0 {
		t.Fatalf("Expected resolved DNS records")
	}
}

func TestDNSProbeRun(t *testing.T) {
	tests := []struct {
		name        string
		probe       DNSProbe
		wantSuccess bool
		timeout     time.Duration
	}{
		{
			name:        "Runs A record probe",
			probe:       dnsProbe("dns.google", DNSRecordTypeA, "8.8.8.8"),
			wantSuccess: true,
		},
		{
			name:        "Runs AAAA record probe",
			probe:       dnsProbe("dns.google", DNSRecordTypeAAAA, "2001:4860:4860::8888"),
			wantSuccess: true,
		},
		{
			name:        "Runs CNAME record probe",
			probe:       dnsProbe("www.github.com", DNSRecordTypeCNAME, "github.com."),
			wantSuccess: true,
		},
		{
			name:        "Runs MX record probe",
			probe:       dnsProbe("gmail.com", DNSRecordTypeMX, "gmail-smtp-in.l.google.com.:5"),
			wantSuccess: true,
		},
		{
			name:        "Runs TXT record probe",
			probe:       dnsProbe("example.com", DNSRecordTypeTXT, "v=spf1 -all"),
			wantSuccess: true,
		},
		{
			name:        "Runs NS record probe",
			probe:       dnsProbe("google.com", DNSRecordTypeNS, "ns1.google.com."),
			wantSuccess: true,
		},
		{
			name:        "Runs SRV record probe",
			probe:       dnsProbe("_minecraft._tcp.hypixel.net", DNSRecordTypeSRV, "mc.hypixel.net.:25565"),
			wantSuccess: true,
		},
		{
			name:        "Fails A record probe with wrong expected value",
			probe:       dnsProbe("dns.google", DNSRecordTypeA, "1.2.3.4"),
			wantSuccess: false,
		},
		{
			name:        "Fails A record probe with unresolvable hostname",
			probe:       dnsProbe("this.hostname.definitely.does.not.exist.invalid", DNSRecordTypeA, "1.2.3.4"),
			wantSuccess: false,
		},
		{
			name: "Fails A record probe with unreachable DNS server",
			probe: DNSProbe{
				Hostname:             "dns.google",
				DNSServer:            "192.0.2.1:53", // non-routable
				RecordType:           DNSRecordTypeA,
				ExpectedRecordValues: []string{"8.8.8.8"},
			},
			wantSuccess: false,
			timeout:     2 * time.Second,
		},
		{
			name:        "Fails AAAA record probe with wrong expected value",
			probe:       dnsProbe("dns.google", DNSRecordTypeAAAA, "::1"),
			wantSuccess: false,
		},
		{
			name:        "Fails CNAME record probe with wrong expected value",
			probe:       dnsProbe("www.github.com", DNSRecordTypeCNAME, "wrong.example.com."),
			wantSuccess: false,
		},
		{
			name:        "Fails CNAME record probe with invalid expected values type",
			probe:       dnsProbe("www.github.com", DNSRecordTypeCNAME, "wrong"),
			wantSuccess: false,
		},
		{
			name:        "Fails MX record probe with wrong expected host",
			probe:       dnsProbe("gmail.com", DNSRecordTypeMX, "nonexistent-mx.example.com."),
			wantSuccess: false,
		},
		{
			name:        "Fails MX record probe with wrong expected priority",
			probe:       dnsProbe("gmail.com", DNSRecordTypeMX, "gmail-smtp-in.l.google.com."),
			wantSuccess: false,
		},
		{
			name:        "Fails TXT record probe with wrong expected value",
			probe:       dnsProbe("example.com", DNSRecordTypeTXT, "this-txt-record-does-not-exist"),
			wantSuccess: false,
		},
		{
			name:        "Fails NS record probe with wrong expected value",
			probe:       dnsProbe("google.com", DNSRecordTypeNS, "ns99.fakens.example.com."),
			wantSuccess: false,
		},
		{
			name:        "Fails SRV record probe with wrong expected target",
			probe:       dnsProbe("_minecraft._tcp.hypixel.net", DNSRecordTypeSRV, "wrong.target.example.com."),
			wantSuccess: false,
		},
		{
			name:        "Fails SRV record probe with wrong expected port",
			probe:       dnsProbe("_minecraft._tcp.hypixel.net", DNSRecordTypeSRV, "mc.hypixel.net."),
			wantSuccess: false,
		},
		{
			name:        "Fails SRV record probe with invalid hostname format",
			probe:       dnsProbe("invalid-srv-hostname", DNSRecordTypeSRV, "something."),
			wantSuccess: false,
		},
		{
			name:        "Fails with unsupported record type",
			probe:       dnsProbe("example.com", DNSRecordType("UNSUPPORTED")),
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDNSProbeRun(t, tt.probe, tt.wantSuccess, tt.timeout)
		})
	}
}

func TestDNSProbeValidate(t *testing.T) {
	tests := []struct {
		name    string
		probe   DNSProbe
		wantErr string
	}{
		{
			name: "Fails when hostname empty",
			probe: DNSProbe{
				Hostname:             "",
				RecordType:           DNSRecordTypeA,
				ExpectedRecordValues: []string{"1.2.3.4"},
			},
			wantErr: "hostname cannot be empty",
		},
		{
			name: "Fails when record type empty",
			probe: DNSProbe{
				Hostname:             "example.com",
				RecordType:           "",
				ExpectedRecordValues: []string{"1.2.3.4"},
			},
			wantErr: "record type cannot be empty",
		},
		{
			name: "Fails when expected values nil",
			probe: DNSProbe{
				Hostname:             "example.com",
				RecordType:           DNSRecordTypeA,
				ExpectedRecordValues: nil,
			},
			wantErr: "expected values must be defined",
		},
		{
			name: "Succeeds with valid config",
			probe: DNSProbe{
				Hostname:             "example.com",
				RecordType:           DNSRecordTypeA,
				ExpectedRecordValues: []string{"1.2.3.4"},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.probe.Validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Expected error")
			}

			if err.Error() != tt.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
