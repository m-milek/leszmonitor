import type { Timestamps } from "@/lib/types";

export const monitorStatuses = ["up", "down", "paused", "maintenance"] as const;
export type MonitorStatus = (typeof monitorStatuses)[number];

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

interface BaseMonitor extends Timestamps {
  id: string;
  name: string;
  slug: string;
  description?: string;
  tagIds: string[];
  ownerId: string;
  interval: number;
  // Retention seconds not configurable yet
  runState: MonitorRunState;
}

export interface HttpMonitor extends BaseMonitor {
  type: "http";
  probeConfig?: HttpMonitorConfig;
}

export interface TcpMonitor extends BaseMonitor {
  type: "tcp";
  probeConfig?: TcpMonitorConfig;
}

export interface DnsMonitor extends BaseMonitor {
  type: "dns";
  probeConfig?: DnsMonitorConfig;
}

export type Monitor = HttpMonitor | TcpMonitor | DnsMonitor;

export type MonitorWithStatus = Monitor & {
  status: MonitorStatus;
};

// Runtime zod schemas and form-value helpers live in
// "@/features/monitors/schema.ts" so that zod is only pulled into the
// route chunks that actually validate monitor forms, keeping it out of the
// initial bundle.
export type {
  MonitorFormValues,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/features/monitors/schema";

export interface HttpResultDetails {
  statusCode: number;
  headers?: Record<string, string>;
  body?: string;
  contentLength: number;
  proto: string;
}

export interface TcpResultDetails {
  latencyMs: number;
}

// results.DNSResultDetails on the server; the records are untyped `any` there.
export interface DnsResultDetails {
  resolvedRecords?: unknown[];
}

export interface MonitorFailure {
  reason: string;
  details?: unknown;
  error?: string;
}

export interface MonitorResult {
  id: string;
  monitorId: string;
  status: MonitorStatus;
  isManuallyTriggered: boolean;
  durationMs: number;
  failures?: MonitorFailure[];
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
  httpCodeToCount: Record<string, number>;
}

export type ProbeSpecificStats = HttpProbeStats;

export interface MonitorStats {
  latency: LatencyStats;
  statusChange: StatusChangeStats;
  uptime: UptimeStats;
  probeSpecific?: ProbeSpecificStats;
  probeType: MonitorType;
}
