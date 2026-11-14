import {
  isValidMonitorType,
  monitorStatuses,
  type MonitorStatus,
  type MonitorType,
} from "@/features/monitors/types";

export interface MonitorsSearch {
  q?: string;
  tags?: string[];
  status?: MonitorStatus[];
  type?: MonitorType[];
}

const parseStringArray = (value: unknown): string[] | undefined => {
  if (!Array.isArray(value)) return undefined;
  const strings = value.filter((v): v is string => typeof v === "string");
  return strings.length > 0 ? strings : undefined;
};

export const parseMonitorsSearch = (
  search: Record<string, unknown>,
): MonitorsSearch => {
  const status = parseStringArray(search.status)?.filter(
    (s): s is MonitorStatus => monitorStatuses.includes(s as MonitorStatus),
  );
  const type = parseStringArray(search.type)?.filter(isValidMonitorType);
  return {
    q: typeof search.q === "string" && search.q !== "" ? search.q : undefined,
    tags: parseStringArray(search.tags),
    status: status && status.length > 0 ? status : undefined,
    type: type && type.length > 0 ? type : undefined,
  };
};
