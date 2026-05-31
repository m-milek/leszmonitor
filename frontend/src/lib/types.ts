import z from "zod";
import { isSlugValid } from "@/lib/slugFromString.ts";

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
  state: MonitorState;
  type: MonitorType;
  probeConfig?: HttpMonitorConfig | TcpMonitorConfig;
}

const monitorStates = ["active", "paused"] as const;
export type MonitorState = (typeof monitorStates)[number];

const monitorTypes = ["http", "tcp", "dns"] as const;
export type MonitorType = (typeof monitorTypes)[number];

export const isValidMonitorType = (value: string): value is MonitorType => {
  return monitorTypes.includes(value as MonitorType);
};

const httpMethods = ["GET", "POST", "PUT", "DELETE", "PATCH"] as const;
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

const tcpProtocols = ["tcp", "tcp4", "tcp6"] as const;
export type TcpProtocol = (typeof tcpProtocols)[number];

export interface TcpMonitorConfig {
  host: string;
  port: number;
  protocol: TcpProtocol;
  timeout: number;
  retryCount: number;
}

export const httpMonitorConfigSchema = z.object({
  method: z.enum(httpMethods),
  url: z.url("Invalid URL"),
  headers: z.record(z.string(), z.string()).optional(),
  body: z.string().optional(),
  saveResponseBody: z.boolean(),
  saveResponseHeaders: z.boolean(),
  expectedStatusCodes: z.array(z.number()).optional(),
  expectedBodyRegex: z.string().optional(),
  expectedHeaders: z.record(z.string(), z.string()).optional(),
  expectedResponseTimeMs: z
    .number()
    .min(1, "Expected response time must be at least 1 ms")
    .optional(),
});

export const tcpMonitorConfigSchema = z.object({
  host: z.string().min(1, "Host is required"),
  port: z
    .number()
    .min(1, "Port must be at least 1")
    .max(65535, "Port must be at most 65535"),
  protocol: z.enum(tcpProtocols),
  timeout: z.number().min(1, "Timeout must be at least 1 ms"),
  retryCount: z.number().min(0, "Retry count cannot be negative"),
});

const recordTypes = ["A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV"] as const;
export type DnsRecordType = (typeof recordTypes)[number];
export const dnsMonitorConfigSchema = z.object({
  hostname: z.string().min(1, "Hostname is required"),
  dnsServer: z.string().min(1, "DNS server address is required"),
  recordType: z.enum(recordTypes),
  expectedRecordValues: z.array(z.string()).default([]),
});

const baseMonitorFields = {
  name: z.string({ message: "Name is required" }).min(1, "Name is required"),
  slug: z
    .string({ message: "Slug is required" })
    .min(1, "Slug is required")
    .refine(
      isSlugValid,
      "Invalid slug format. Must be lowercase, alphanumeric, and can include hyphens.",
    ),
  description: z.string().optional(),
  projectSlug: z.string(),
  interval: z
    .number({ message: "Interval must be a number" })
    .min(1, "Interval must be at least 1 second"),
};

const httpMonitorSchema = z.object({
  ...baseMonitorFields,
  type: z.literal("http"),
  probeConfig: httpMonitorConfigSchema.optional(),
});

const tcpMonitorSchema = z.object({
  ...baseMonitorFields,
  type: z.literal("tcp"),
  probeConfig: tcpMonitorConfigSchema.optional(),
});

const dnsMonitorSchema = z.object({
  ...baseMonitorFields,
  type: z.literal("dns"),
  probeConfig: dnsMonitorConfigSchema.optional(),
});

export const newMonitorSchema = z.discriminatedUnion("type", [
  httpMonitorSchema,
  tcpMonitorSchema,
  dnsMonitorSchema,
]);

export type MonitorFormValues = z.infer<typeof newMonitorSchema>;

export const newMonitorSchemaDefaultValues = {
  name: "",
  slug: "",
  description: "",
  projectSlug: "",
  interval: 60,
} satisfies Partial<MonitorFormValues>;

export const defaultConfigs: Record<
  MonitorType,
  MonitorFormValues["probeConfig"]
> = {
  http: {
    method: "GET",
    url: "",
    headers: {},
    body: "",
    saveResponseBody: false,
    saveResponseHeaders: false,
    expectedStatusCodes: [],
    expectedHeaders: {},
  },
  tcp: {
    host: "",
    port: 443,
    protocol: "tcp",
    timeout: 5000,
    retryCount: 3,
  },
  dns: {
    hostname: "",
    recordType: "A",
    dnsServer: "1.1.1.1",
    expectedRecordValues: [],
  },
};

export type MonitorCreatePayload = MonitorFormValues;
export type MonitorUpdatePayload = MonitorFormValues & { id: string };

export const mapMonitorToFormValues = (monitor: Monitor): MonitorFormValues => {
  const configDefaults = defaultConfigs[monitor.type];

  return {
    ...newMonitorSchemaDefaultValues,
    projectSlug: monitor.projectSlug,
    name: monitor.name,
    slug: monitor.slug,
    description: monitor.description ?? "",
    interval: monitor.interval,
    type: monitor.type,
    probeConfig: {
      ...configDefaults,
      ...monitor.probeConfig,
    },
  } as MonitorFormValues;
};

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
  isSuccess: boolean;
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
