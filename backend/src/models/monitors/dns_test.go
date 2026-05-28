package monitors

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

func TestDNSProbeRun(t *testing.T) {
	t.Run("Runs A record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeA,
			ExpectedValues: DNSAExpectedValues{"8.8.8.8"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs AAAA record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeAAAA,
			ExpectedValues: DNSAAAAExpectedValues{"2001:4860:4860::8888"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs CNAME record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "www.github.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeCNAME,
			ExpectedValues: DNSCNAMEExpectedValues{CNAME: "github.com."},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs MX record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "gmail.com",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeMX,
			ExpectedValues: DNSMXExpectedValues{
				{Host: "gmail-smtp-in.l.google.com.", Priority: 5},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs TXT record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "example.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeTXT,
			ExpectedValues: DNSTXTExpectedValues{"v=spf1 -all"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs NS record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "google.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeNS,
			ExpectedValues: DNSNSExpectedValues{"ns1.google.com."},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Runs SRV record probe", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "_minecraft._tcp.hypixel.net",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeSRV,
			ExpectedValues: DNSSRVExpectedValues{
				{Target: "mc.hypixel.net.", Port: 25565, Priority: 0, Weight: 0},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if !result.GetIsSuccess() {
			t.Fatalf("Expected success, got errors: %+v", result.GetErrorDetails())
		}

		details, ok := result.GetDetails().(*monitorresult.DNSResultDetails)
		if !ok {
			t.Fatalf("Expected DNS result details")
		}
		if len(details.ResolvedRecords) == 0 {
			t.Fatalf("Expected resolved DNS records")
		}
	})

	t.Run("Fails A record probe with wrong expected value", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeA,
			ExpectedValues: DNSAExpectedValues{"1.2.3.4"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails A record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeA,
			ExpectedValues: "not-a-slice",
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails A record probe with unresolvable hostname", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "this.hostname.definitely.does.not.exist.invalid",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeA,
			ExpectedValues: DNSAExpectedValues{"1.2.3.4"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails A record probe with unreachable DNS server", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "192.0.2.1:53", // non-routable
			RecordType:     DNSRecordTypeA,
			ExpectedValues: DNSAExpectedValues{"8.8.8.8"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails AAAA record probe with wrong expected value", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeAAAA,
			ExpectedValues: DNSAAAAExpectedValues{"::1"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails AAAA record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "dns.google",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeAAAA,
			ExpectedValues: 12345,
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails CNAME record probe with wrong expected value", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "www.github.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeCNAME,
			ExpectedValues: DNSCNAMEExpectedValues{CNAME: "wrong.example.com."},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails CNAME record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "www.github.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeCNAME,
			ExpectedValues: []string{"wrong"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails MX record probe with wrong expected host", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "gmail.com",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeMX,
			ExpectedValues: DNSMXExpectedValues{
				{Host: "nonexistent-mx.example.com.", Priority: 5},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails MX record probe with wrong expected priority", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "gmail.com",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeMX,
			ExpectedValues: DNSMXExpectedValues{
				{Host: "gmail-smtp-in.l.google.com.", Priority: 999},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails MX record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "gmail.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeMX,
			ExpectedValues: "not-valid",
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails TXT record probe with wrong expected value", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "example.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeTXT,
			ExpectedValues: DNSTXTExpectedValues{"this-txt-record-does-not-exist"},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails TXT record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "example.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeTXT,
			ExpectedValues: 42,
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails NS record probe with wrong expected value", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "google.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeNS,
			ExpectedValues: DNSNSExpectedValues{"ns99.fakens.example.com."},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails NS record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "google.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeNS,
			ExpectedValues: true,
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails SRV record probe with wrong expected target", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "_minecraft._tcp.hypixel.net",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeSRV,
			ExpectedValues: DNSSRVExpectedValues{
				{Target: "wrong.target.example.com.", Port: 25565},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails SRV record probe with wrong expected port", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "_minecraft._tcp.hypixel.net",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeSRV,
			ExpectedValues: DNSSRVExpectedValues{
				{Target: "mc.hypixel.net.", Port: 9999},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails SRV record probe with invalid expected values type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "_minecraft._tcp.hypixel.net",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordTypeSRV,
			ExpectedValues: "invalid",
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails SRV record probe with invalid hostname format", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:   "invalid-srv-hostname",
			DNSServer:  "8.8.8.8:53",
			RecordType: DNSRecordTypeSRV,
			ExpectedValues: DNSSRVExpectedValues{
				{Target: "something.", Port: 80},
			},
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})

	t.Run("Fails with unsupported record type", func(t *testing.T) {
		probe := DNSProbe{
			Hostname:       "example.com",
			DNSServer:      "8.8.8.8:53",
			RecordType:     DNSRecordType("UNSUPPORTED"),
			ExpectedValues: nil,
		}

		result := probe.Run(t.Context(), uuid.New())
		if result.GetIsSuccess() {
			t.Fatalf("Expected failure, got success")
		}
	})
}
