import { z } from "zod";
import { monitorStatuses, monitorTypes } from "@/features/monitors/types";
import { searchArray } from "@/lib/search-params";

export const monitorsSearchSchema = z.object({
  q: z.string().optional().catch(undefined),
  tags: searchArray(z.string()),
  status: searchArray(z.enum(monitorStatuses)),
  type: searchArray(z.enum(monitorTypes)),
});

export type MonitorsSearch = z.infer<typeof monitorsSearchSchema>;
