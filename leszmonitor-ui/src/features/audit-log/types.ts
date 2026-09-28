export interface AuditLogEntry {
  id: string;
  username?: string;
  resourceId?: string;
  action: string;
  isSuccess: boolean;
  summary?: string;
  before?: string;
  after?: string;
  traceId?: string;
  createdAt: Date;
}

export interface AuditLogFilters {
  username?: string;
  resourceId?: string;
  action?: string;
  isSuccess?: boolean;
  startDate?: Date;
  endDate?: Date;
}

export const auditLogActions = [
  "monitor.create",
  "monitor.update",
  "monitor.delete",
  "monitor.run_manual",
  "user.create",
  "user.update",
  "user.delete",
  "tag.create",
  "tag.update",
  "tag.delete",
  "auth.login",
  "auth.failed_login",
  "auth.password_change",
  "global_parameter.update",
] as const;
