
export interface Timestamps {
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends Timestamps {
  id: string;
  username: string;
}

export enum ProjectRole {
  Owner = "owner",
  Admin = "admin",
  Member = "member",
  Viewer = "viewer",
}

export const mapProjectRoleToDisplayName: Record<ProjectRole, string> = {
  [ProjectRole.Owner]: "Owner",
  [ProjectRole.Admin]: "Admin",
  [ProjectRole.Member]: "Member",
  [ProjectRole.Viewer]: "Viewer",
};

export interface ProjectMember extends Timestamps {
  id: string;
  username: string;
  role: ProjectRole;
}

export interface Project extends Timestamps {
  id: string;
  slug: string;
  name: string;
  description: string;
  members: ProjectMember[];
}

export interface Monitor extends Timestamps {
  id: string;
  name: string;
  slug: string;
  description?: string;
  projectSlug: string;
  interval: number;
  // Retention seconds not configurable yet
  runState: MonitorRunState;
  type: MonitorType;
  probeConfig?: HttpMonitorConfig | TcpMonitorConfig;
}

const monitorStatuses = ["up", "down", "paused", "maintenance"] as const;
export type MonitorStatus = (typeof monitorStatuses)[number];

const monitorRunStates = ["active", "paused"] as const;
export type MonitorRunState = (typeof monitorRunStates)[number];

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

// Runtime zod schemas and form-value helpers live in "@/lib/monitorSchema.ts"
// so that zod is only pulled into the route chunks that actually validate
// monitor forms, keeping it out of the initial bundle.
export type {
  MonitorFormValues,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/lib/monitorSchema.ts";

export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  jwt: string;
}

export interface JwtClaims {
  username: string;
}

export interface MonitorResult {
  id: string;
  monitorId: string;
  status: MonitorStatus;
  isManuallyTriggered: boolean;
  durationMs: number;
  errorDetails: ErrorDetails;
  monitorType: string;
  details: HttpResultDetails | TcpResultDetails;
  createdAt: Date;
}

export interface ErrorDetails {
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
    typeof (obj as any).type === "string" &&
    "monitorId" in obj &&
    typeof (obj as any).monitorId === "string" &&
    "response" in obj &&
    typeof (obj as any).response === "object"
  );
};

export interface Pagination {
  page: number;
  perPage: number;
}

export interface ErrorDetails {
  message: string;
}

export interface ApiError {
  error: ErrorDetails;
  status: number;
}

export interface AuditLogEntry {
  id: string;
  username?: string;
  projectId?: string;
  resourceId?: string;
  action: string;
  isSuccess: boolean;
  summary?: string;
  before?: string;
  after?: string;
  traceId?: string;
  createdAt: Date;
}

export interface AuditLogFilters {
  username?: string;
  projectId?: string;
  resourceId?: string;
  action?: string;
  isSuccess?: boolean;
  startDate?: Date;
  endDate?: Date;
}
