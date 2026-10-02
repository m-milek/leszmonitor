import type { HttpMonitor, HttpMonitorConfig } from "@/features/monitors/types";
import { ConfigTable } from "@/features/monitors/components/DetailsPage/ConfigTable";
import {
  type FieldConfig,
  formatBoolean,
  formatList,
  formatMs,
} from "@/features/monitors/components/DetailsPage/config-fields";

export interface HttpConfigContentProps {
  monitor: HttpMonitor;
}

const renderKeyValues = (record: Record<string, string>) => (
  <ul>
    {Object.entries(record).map(([key, value]) => (
      <li key={key}>
        <span className="font-medium">{key}:</span> {value}
      </li>
    ))}
  </ul>
);

const httpFields: FieldConfig<HttpMonitorConfig> = {
  url: {
    label: "Endpoint",
    render: (url, config) => (
      <pre className="flex gap-2">
        <span>{config.method}</span>
        <span>{url}</span>
      </pre>
    ),
  },
  headers: { label: "Headers", render: renderKeyValues },
  body: {
    label: "Body",
    render: (value) => <pre className="font-mono text-xs">{value}</pre>,
  },
  saveResponseBody: { label: "Save Response Body", render: formatBoolean },
  saveResponseHeaders: {
    label: "Save Response Headers",
    render: formatBoolean,
  },
  expectedStatusCodes: {
    label: "Expected Status Codes",
    render: formatList,
  },
  expectedBodyRegex: {
    label: "Expected Body Pattern",
    render: (value) => <code className="font-mono">{value}</code>,
  },
  expectedHeaders: { label: "Expected Headers", render: renderKeyValues },
  expectedResponseTimeMs: {
    label: "Expected Response Time",
    render: formatMs,
  },
};

export const HttpConfigContent = ({ monitor }: HttpConfigContentProps) => {
  if (!monitor.probeConfig) return null;
  return <ConfigTable data={monitor.probeConfig} fields={httpFields} />;
};
