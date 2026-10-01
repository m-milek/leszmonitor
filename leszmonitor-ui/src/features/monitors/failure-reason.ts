import type { FailureReason } from "@/features/monitors/types";

const FAILURE_REASON_LABELS: Record<FailureReason, string> = {
  HTTP_REQUEST_FAILED: "Request failed",
  HTTP_RESPONSE_BODY_READ_FAILED: "Could not read response body",
  HTTP_STATUS_CODE_MISMATCH: "Unexpected status code",
  HTTP_RESPONSE_BODY_MISMATCH: "Unexpected response body",
  HTTP_RESPONSE_HEADER_MISMATCH: "Unexpected response headers",
  HTTP_RESPONSE_TIME_EXCEEDED: "Response too slow",
  DNS_LOOKUP_FAILED: "DNS lookup failed",
  DNS_INVALID_SRV_HOSTNAME: "Invalid SRV hostname",
  DNS_EXPECTED_RECORD_MISSING: "Expected record missing",
  TCP_CONNECTION_FAILED: "Connection failed",
};

export const failureReasonToLabel = (reason: FailureReason): string =>
  FAILURE_REASON_LABELS[reason] ?? reason;
