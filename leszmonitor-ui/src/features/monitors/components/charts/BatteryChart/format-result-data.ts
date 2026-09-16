import type { MonitorResult } from "@/features/monitors/types";

export const formatResultData = (
  result: MonitorResult,
): Record<string, string> => {
  const isError = result.status === "down" && result.errorDetails;
  return {
    ID: result.id,
    Status: result.status,
    "Created At": result.createdAt.toLocaleString(),
    Details: isError
      ? JSON.stringify(result.errorDetails, null, 2)
      : JSON.stringify(result.details, null, 2),
  };
};
