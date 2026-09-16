import { AuditLogApi } from "@/features/audit-log/audit-log-api";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { AuditLogTable } from "@/features/audit-log/components/AuditLogTable";
import type { AuditLogFilters } from "@/features/audit-log/types";

export function AuditLogPage() {
  const [filters] = useState<AuditLogFilters>({});

  const { data: logs } = useQuery({
    queryKey: ["auditLogs", filters],
    queryFn: () => AuditLogApi.getByFilter(filters),
  });

  return (
    <PageContainer>
      <TypographyH1>Audit Log</TypographyH1>
      <AuditLogTable entries={logs ?? []} />
    </PageContainer>
  );
}
