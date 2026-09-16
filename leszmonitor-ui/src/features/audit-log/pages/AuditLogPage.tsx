import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1 } from "@/components/common/Typography.tsx";
import { AuditLogTable } from "@/features/audit-log/components/AuditLogTable.tsx";
import { getAuditLogByFilter } from "@/features/audit-log/audit-log-api.ts";
import type { AuditLogFilters } from "@/features/audit-log/types.ts";

export function AuditLogPage() {
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
