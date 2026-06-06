import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState, useMemo } from "react";
import { getAuditLogByFilter } from "@/lib/data/auditLogData.ts";
import type { AuditLogFilters } from "@/lib/types.ts";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { AuditLogTable } from "@/components/leszmonitor/tables/audit-log/AuditLogTable.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/audit-log",
)({
  component: RouteComponent,
});

const parentRoute = getRouteApi("/_authenticated/projects/$projectId");

function RouteComponent() {
  const project = parentRoute.useLoaderData();

  const [extraFilters] = useState<Omit<AuditLogFilters, "projectId">>({});

  const filters = useMemo<AuditLogFilters>(
    () => ({ ...extraFilters, projectId: project.id }),
    [extraFilters, project.id],
  );

  const { data: logs } = useQuery({
    queryKey: ["auditLogs", project.id, extraFilters, filters],
    queryFn: () => getAuditLogByFilter(filters),
  });

  return (
    <MainPanelContainer>
      <pre>{JSON.stringify(logs, null, 2)}</pre>
      <AuditLogTable entries={logs ?? []} />
    </MainPanelContainer>
  );
}
