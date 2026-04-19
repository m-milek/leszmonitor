import z from "zod";

export interface Timestamps {
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends Timestamps {
  id: string;
  username: string;
}

export enum OrgRole {
  Owner = "owner",
  Admin = "admin",
  Member = "member",
  Viewer = "viewer",
}

export const mapOrgRoleToDisplayName: Record<OrgRole, string> = {
  [OrgRole.Owner]: "Owner",
  [OrgRole.Admin]: "Admin",
  [OrgRole.Member]: "Member",
  [OrgRole.Viewer]: "Viewer",
};

export interface OrgMember extends Timestamps {
  id: string;
  username: string;
  role: OrgRole;
}

export interface Org extends Timestamps {
  id: string;
  displayId: string;
  name: string;
  description: string;
  members: OrgMember[];
}

export interface Project extends Timestamps {
  id: string;
  name: string;
  displayId: string;
  description?: string;
}

export interface Monitor extends Timestamps {
  id: string;
  name: string;
  displayId: string;
  description?: string;
  projectId?: string;
  orgId: string;
  interval: number;
  type: MonitorType;
  config?: HttpMonitorConfig | PingMonitorConfig;
}

export type MonitorType = "http" | "ping";

export const isValidMonitorType = (value: string): value is MonitorType => {
  const values = ["http", "ping"] as MonitorType[];
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

export interface PingMonitorConfig {
  host: string;
  port: number;
  protocol: "tcp" | "udp" | "tcp4" | "tcp6" | "udp4" | "udp6";
  timeoutMs: number;
  retryCount: number;
}

export const httpMonitorConfigSchema = z.object({
  method: z.enum(["GET", "POST", "PUT", "DELETE", "PATCH"]),
  url: z.url("Invalid URL"),
  headers: z.record(z.string(), z.string()).optional(),
  body: z.string().optional(),
  saveResponseBody: z.boolean().default(false),
  saveResponseHeaders: z.boolean().default(false),
  expectedStatusCodes: z.array(z.number()).default([]),
  expectedBodyRegex: z.string().optional(),
  expectedHeaders: z.record(z.string(), z.string()).default({}),
  expectedResponseTimeMs: z
    .number()
    .min(1, "Expected response time must be at least 1 ms")
    .optional(),
});

export const pingMonitorConfigSchema = z.object({
  host: z.string().min(1, "Host is required"),
  port: z
    .number()
    .min(1, "Port must be at least 1")
    .max(65535, "Port must be at most 65535"),
  protocol: z.enum(["tcp", "udp", "tcp4", "tcp6", "udp4", "udp6"]),
  timeoutMs: z.number().min(1, "Timeout must be at least 1 ms"),
  retryCount: z.number().min(0, "Retry count cannot be negative"),
});

const baseMonitorFields = {
  name: z.string({ message: "Name is required" }).min(1, "Name is required"),
  displayId: z
    .string({ message: "Display ID is required" })
    .min(1, "Display ID is required"),
  description: z.string().optional(),
  orgId: z.string(),
  projectId: z.string(),
  interval: z
    .number({ message: "Interval must be a number" })
    .min(1, "Interval must be at least 1 second"),
};

const httpMonitorSchema = z.object({
  ...baseMonitorFields,
  type: z.literal("http"),
  config: httpMonitorConfigSchema.optional(),
});

const pingMonitorSchema = z.object({
  ...baseMonitorFields,
  type: z.literal("ping"),
  config: pingMonitorConfigSchema.optional(),
});

export const newMonitorSchema = z.discriminatedUnion("type", [
  httpMonitorSchema,
  pingMonitorSchema,
]);

export type MonitorFormValues = z.infer<typeof newMonitorSchema>;
export type HttpMonitorFormValues = z.infer<typeof httpMonitorSchema>;
export type PingMonitorFormValues = z.infer<typeof pingMonitorSchema>;

export const newMonitorSchemaDefaultValues = {
  name: "",
  displayId: "",
  description: "",
  orgId: "",
  projectId: "",
  interval: 60,
} satisfies Partial<MonitorFormValues>;

export const defaultConfigs: Record<MonitorType, MonitorFormValues["config"]> =
  {
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
    ping: {
      host: "",
      port: 443,
      protocol: "tcp",
      timeoutMs: 5000,
      retryCount: 3,
    },
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
