import type { MonitorStatus } from "@/features/monitors/types";
import type { Tag } from "@/features/tags/types";
import type { ResourceSource, Timestamps } from "@/lib/types";

export interface StatusPage extends Timestamps {
  id: string;
  slug: string;
  title: string;
  description: string;
  groupMonitorsByTag: boolean;
  requireAuth: boolean;
  enabled: boolean;
  source: ResourceSource;
  tagIds: string[];
  monitorCount: number;
}

export interface PublicMonitorResult {
  status: MonitorStatus;
  createdAt: Date;
}

export interface PublicMonitor {
  id: string;
  name: string;
  status: MonitorStatus;
  uptimePercent: number | null;
  tags: Tag[];
  results: PublicMonitorResult[];
}
