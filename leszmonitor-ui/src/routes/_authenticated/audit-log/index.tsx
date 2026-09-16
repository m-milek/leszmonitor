import { createFileRoute } from "@tanstack/react-router";
import { AuditLogPage } from "@/features/audit-log/pages/AuditLogPage.tsx";

export const Route = createFileRoute("/_authenticated/audit-log/")({
  head: () => ({
    meta: [{ title: "Audit Log | Leszmonitor" }],
  }),
  component: AuditLogPage,
});
