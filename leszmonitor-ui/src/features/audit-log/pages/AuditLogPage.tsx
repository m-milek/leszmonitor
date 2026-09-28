import { AuditLogApi } from "@/features/audit-log/audit-log-api";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { LucideX } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { AuditLogTable } from "@/features/audit-log/components/AuditLogTable";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { LMSelect } from "@/components/form/LMSelect";
import { UsersApi } from "@/features/users/users-api";
import { QUERY_KEYS } from "@/lib/consts";
import {
  auditLogActions,
  type AuditLogFilters,
} from "@/features/audit-log/types";

const statusOptions = [
  { value: "true", label: "Success" },
  { value: "false", label: "Failed" },
];

export function AuditLogPage() {
  const [filters, setFilters] = useState<AuditLogFilters>({});

  const { data: logs } = useQuery({
    queryKey: ["auditLogs", filters],
    queryFn: () => AuditLogApi.getByFilter(filters),
  });

  const { data: users = [] } = useQuery({
    queryKey: [QUERY_KEYS.USERS],
    queryFn: () => UsersApi.getAll(),
  });

  const hasActiveFilters = Object.values(filters).some(
    (value) => value !== undefined,
  );

  return (
    <PageContainer>
      <TypographyH1>Audit Log</TypographyH1>
      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center gap-2">
            <LMSelect
              id="audit-log-user-filter"
              name="audit-log-user-filter"
              placeholder="User"
              className="w-48"
              items={users.map((user) => ({
                value: user.username,
                label: user.username,
              }))}
              value={filters.username ?? ""}
              onValueChange={(username) =>
                setFilters({ ...filters, username: username || undefined })
              }
            />
            <LMSelect
              id="audit-log-action-filter"
              name="audit-log-action-filter"
              placeholder="Action"
              className="w-56"
              items={auditLogActions.map((action) => ({
                value: action,
                label: action,
              }))}
              value={filters.action ?? ""}
              onValueChange={(action) =>
                setFilters({ ...filters, action: action || undefined })
              }
            />
            <LMSelect
              id="audit-log-status-filter"
              name="audit-log-status-filter"
              placeholder="Status"
              className="w-36"
              items={statusOptions}
              value={filters.isSuccess?.toString() ?? ""}
              onValueChange={(isSuccess) =>
                setFilters({
                  ...filters,
                  isSuccess: isSuccess ? isSuccess === "true" : undefined,
                })
              }
            />
            {hasActiveFilters && (
              <Button variant="ghost" onClick={() => setFilters({})}>
                Clear
                <LucideX />
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <AuditLogTable entries={logs ?? []} />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
