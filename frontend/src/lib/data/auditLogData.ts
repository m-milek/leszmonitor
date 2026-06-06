import type { AuditLogEntry, AuditLogFilters } from "@/lib/types.ts";
import { authFetch } from "@/lib/data/utils.ts";
import { BACKEND_API_URL } from "@/lib/consts.ts";

const filterIntoParams = (filter: AuditLogFilters): URLSearchParams => {
  const params = new URLSearchParams();
  Object.entries(filter).forEach(([key, value]) => {
    if (value) {
      if (value instanceof Date) {
        params.append(key, value.toISOString());
      } else {
        params.append(key, value);
      }
    }
  });
  return params;
};

export const getAuditLogByFilter = async (
  filter: AuditLogFilters,
): Promise<AuditLogEntry[]> => {
  const queryParams = filterIntoParams(filter);

  const response = await authFetch(
    `${BACKEND_API_URL}/audit-log?${queryParams.toString()}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );
  if (!response.ok) {
    throw new Error("Failed to fetch audit logs");
  }

  const data: AuditLogEntry[] = await response.json();

  return data.map((entry) => ({
    ...entry,
    createdAt: new Date(entry.createdAt),
  }));
};
