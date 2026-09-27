import { Badge } from "@/components/ui/badge";
import type { GlobalParameterSource } from "@/features/instance/global-parameters-api";

export interface GlobalParameterSourceBadgeProps {
  source: GlobalParameterSource;
}

export function GlobalParameterSourceBadge({
  source,
}: Readonly<GlobalParameterSourceBadgeProps>) {
  if (source === "config") {
    return <Badge variant="outline">Set in config file</Badge>;
  }
  if (source === "default") {
    return <Badge variant="secondary">Default</Badge>;
  }
  return null;
}
