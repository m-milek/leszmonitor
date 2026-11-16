import type { ResourceSource, Timestamps } from "@/lib/types";

export interface StatusPage extends Timestamps {
  id: string;
  slug: string;
  title: string;
  description: string;
  groupMonitorsByTag: boolean;
  requireAuth: boolean;
  source: ResourceSource;
  tagIds: string[];
}
