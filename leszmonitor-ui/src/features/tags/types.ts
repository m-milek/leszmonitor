import type { ResourceSource, Timestamps } from "@/lib/types";

export interface Tag extends Timestamps {
  id: string;
  name: string;
  description: string;
  colorHex: string;
  source: ResourceSource;
}
