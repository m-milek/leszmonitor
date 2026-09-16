import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { getAuditLogByFilter } from "@/lib/data/auditLogData.ts";
import type { AuditLogFilters } from "@/lib/types.ts";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { AuditLogTable } from "@/components/leszmonitor/tables/audit-log/AuditLogTable.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";

export const Route = createFileRoute("/_authenticated/audit-log/")({
  head: () => ({
    meta: [{ title: "Audit Log | Leszmonitor" }],
  }),
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
      <TypographyH1>Audit Log</TypographyH1>
      <AuditLogTable entries={logs ?? []} />
    </PageContainer>
  );
}
