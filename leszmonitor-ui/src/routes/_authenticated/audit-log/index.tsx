import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { getAuditLogByFilter } from "@/lib/data/auditLogData.ts";
import type { AuditLogFilters } from "@/lib/types.ts";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { AuditLogTable } from "@/components/leszmonitor/tables/audit-log/AuditLogTable.tsx";

export const Route = createFileRoute("/_authenticated/audit-log/")({
  component: RouteComponent,
});

function RouteComponent() {
  const [filters] = useState<AuditLogFilters>({});

  const { data: logs } = useQuery({
    queryKey: ["auditLogs", filters],
    queryFn: () => getAuditLogByFilter(filters),
  });

  return (
    <PageContainer>
      <AuditLogTable entries={logs ?? []} />
    </PageContainer>
  );
}
