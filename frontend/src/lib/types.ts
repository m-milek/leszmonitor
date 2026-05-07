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
  type: MonitorType;
  probeConfig?: HttpMonitorConfig | TcpMonitorConfig;
}

export type MonitorType = "http" | "tcp";

export const isValidMonitorType = (value: string): value is MonitorType => {
  const values = ["http", "tcp"];
  return values.includes(value as MonitorType);
};

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

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

export interface TcpMonitorConfig {
  host: string;
  port: number;
  protocol: "tcp" | "tcp4" | "tcp6";
  timeout: number;
  retryCount: number;
}

export const httpMonitorConfigSchema = z.object({
  method: z.enum(["GET", "POST", "PUT", "DELETE", "PATCH"]),
  url: z.url("Invalid URL"),
  headers: z.record(z.string(), z.string()).optional(),
  body: z.string().optional(),
  saveResponseBody: z.boolean().default(false),
  saveResponseHeaders: z.boolean().default(false),
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
  protocol: z.enum(["tcp", "tcp4", "tcp6"]),
  timeout: z.number().min(1, "Timeout must be at least 1 ms"),
  retryCount: z.number().min(0, "Retry count cannot be negative"),
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

export const newMonitorSchema = z.discriminatedUnion("type", [
  httpMonitorSchema,
  tcpMonitorSchema,
]);

export type MonitorFormValues = z.infer<typeof newMonitorSchema>;
export type HttpMonitorFormValues = z.infer<typeof httpMonitorSchema>;
export type TcpMonitorFormValues = z.infer<typeof tcpMonitorSchema>;

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
      ...(monitor.probeConfig ?? {}),
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
