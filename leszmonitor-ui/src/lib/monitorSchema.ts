import z from "zod";
import { isSlugValid } from "@/lib/slugFromString.ts";
import {
  httpMethods,
  tcpProtocols,
  recordTypes,
  type Monitor,
  type MonitorType,
} from "@/lib/types.ts";

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
  resultRetentionSeconds: z
    .number({ message: "Retention period must be a number" })
    .min(1, "Retention period must be at least 1 second")
    .optional(),
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
  resultRetentionSeconds: 43200,
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
