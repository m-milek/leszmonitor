import type { Timestamps } from "@/lib/types";

export type TagSource = "ui" | "config";

export interface Tag extends Timestamps {
  id: string;
  name: string;
  description: string;
  colorHex: string;
  source: TagSource;
}
