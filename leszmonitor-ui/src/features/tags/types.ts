import type { Timestamps } from "@/lib/types.ts";

export interface Tag extends Timestamps {
  id: string;
  name: string;
  description: string;
  colorHex: string;
}
