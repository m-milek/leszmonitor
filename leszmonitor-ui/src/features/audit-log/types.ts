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
