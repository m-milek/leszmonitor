import type { Timestamps } from "@/lib/types";

export type MonitorStatus = "up" | "down" | "paused" | "maintenance";

export type MonitorRunState = "active" | "paused";

const monitorTypes = ["http", "tcp", "dns"] as const;
export type MonitorType = (typeof monitorTypes)[number];

export const isValidMonitorType = (value: string): value is MonitorType => {
  return monitorTypes.includes(value as MonitorType);
};

export const httpMethods = ["GET", "POST", "PUT", "DELETE", "PATCH"] as const;
export type HttpMethod = (typeof httpMethods)[number];

export interface HttpMonitorConfig {
  method: HttpMethod;
  url: string;
  headers?: Record<string, string>;
  body?: string;
  saveResponseBody?: boolean;
  saveResponseHeaders?: boolean;
  expectedStatusCodes?: number[];
  expectedBodyRegex?: string;
  expectedHeaders?: Record<string, string>;
  expectedResponseTimeMs?: number;
}

export const tcpProtocols = ["tcp", "tcp4", "tcp6"] as const;
export type TcpProtocol = (typeof tcpProtocols)[number];

export interface TcpMonitorConfig {
  host: string;
  port: number;
  protocol: TcpProtocol;
  timeout: number;
  retryCount: number;
}

export const recordTypes = [
  "A",
  "AAAA",
  "CNAME",
  "MX",
  "TXT",
  "NS",
  "SRV",
] as const;
export type DnsRecordType = (typeof recordTypes)[number];

export interface DnsMonitorConfig {
  hostname: string;
  dnsServer?: string;
  recordType: DnsRecordType;
  expectedRecordValues: string[];
}

export interface Monitor extends Timestamps {
  id: string;
  name: string;
  slug: string;
  description?: string;
  tagIds?: string[];
  ownerId: string;
  interval: number;
  // Retention seconds not configurable yet
  runState: MonitorRunState;
  type: MonitorType;
  probeConfig?: HttpMonitorConfig | TcpMonitorConfig | DnsMonitorConfig;
}

// Runtime zod schemas and form-value helpers live in
// "@/features/monitors/schema.ts" so that zod is only pulled into the
// route chunks that actually validate monitor forms, keeping it out of the
// initial bundle.
export type {
  MonitorFormValues,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/features/monitors/schema";

export interface MonitorErrorDetails {
  errorMessage: string;
  errors: string[];
  failures: string[];
}

export interface HttpResultDetails {
  statusCode: number;
  headers?: Record<string, string>;
  body?: string;
  contentLength: number;
  proto: string;
}

export interface TcpResultDetails {
  tries: number;
  latencyMs: number;
}

// results.DNSResultDetails on the server; the records are untyped `any` there.
export interface DnsResultDetails {
  resolvedRecords?: unknown[];
}

export interface MonitorResult {
  id: string;
  monitorId: string;
  status: MonitorStatus;
  isManuallyTriggered: boolean;
  durationMs: number;
  errorDetails: MonitorErrorDetails;
  monitorType: string;
  details: HttpResultDetails | TcpResultDetails | DnsResultDetails;
  createdAt: Date;
}

export interface MonitorResultMessage {
  type: string;
  monitorId: string;
  response: MonitorResult;
}

export const isMonitorResultMessage = (
  obj: object,
): obj is MonitorResultMessage => {
  return (
    typeof obj === "object" &&
    obj !== null &&
    "type" in obj &&
    typeof obj.type === "string" &&
    "monitorId" in obj &&
    typeof obj.monitorId === "string" &&
    "response" in obj &&
    typeof obj.response === "object"
  );
};

export interface LatencyStats {
  avg: number;
  min: number;
  max: number;
}
export interface StatusChangeStats {
  secondsInCurrentStatus: number;
}
export interface UptimeStats {
  statusToCount: Partial<Record<MonitorStatus, number>> | null;
  statusToPercentage: Partial<Record<MonitorStatus, number>> | null;
}

export interface HttpProbeStats {
  type: "http";
  httpCodeToCount: Record<string, number>;
}

export type ProbeSpecificStats = HttpProbeStats;

export interface MonitorStats {
  latency: LatencyStats;
  statusChange: StatusChangeStats;
  uptime: UptimeStats;
  probeSpecific?: ProbeSpecificStats;
}
