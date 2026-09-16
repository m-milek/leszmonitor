import { monitorStatusToStatusDot } from "@/features/monitors/model/status.ts";
import type { Monitor } from "@/features/monitors/model/types.ts";
import { TypographyH3 } from "@/components/common/Typography.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { StyledLink } from "@/components/common/StyledLink.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { LucideEdit, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { StatusDot } from "@/components/common/StatusDot.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getLatestMonitorResultByMonitorId } from "@/features/monitors/api/results.ts";
import { useQuery } from "@tanstack/react-query";
import { MonitorStatusPill } from "@/features/monitors/components/MonitorStatusPill.tsx";

export interface MonitorListItemProps {
  monitor: Monitor;
  onDeleteMonitor?: (monitorId: string) => Promise<void>;
  navigateToEditMonitor?: (monitorId: string) => void;
}

export function MonitorListItem({
  monitor,
  onDeleteMonitor,
  navigateToEditMonitor,
}: Readonly<MonitorListItemProps>) {
  const { data: lastResultData } = useQuery({
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id],
    queryFn: () => getLatestMonitorResultByMonitorId(monitor.id),
  });

  const dotStatus = monitorStatusToStatusDot(lastResultData?.status);

  return (
    <Card>
      <CardHeader>
        <Flex direction="row" className="justify-between">
          <Flex direction="row" className="items-center gap-2">
            <StatusDot status={dotStatus} />
            <TypographyH3>
              <StyledLink
                to="/monitors/$monitorSlug"
                params={{ monitorSlug: monitor.slug }}
              >
                {monitor.name}
              </StyledLink>
            </TypographyH3>
            <MonitorStatusPill monitor={monitor} />
          </Flex>
          <Flex direction="row">
            {navigateToEditMonitor && onDeleteMonitor && (
              <>
                <Button
                  variant="ghost"
                  size="icon-lg"
                  onClick={() => navigateToEditMonitor(monitor.slug)}
                >
                  <LucideEdit className="size-5" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon-lg"
                  onClick={() => onDeleteMonitor(monitor.id)}
                >
                  <LucideTrash2 className="size-5 text-destructive" />
                </Button>
              </>
            )}
          </Flex>
        </Flex>
      </CardHeader>
      <CardContent>
        <Flex direction="column">
          <span>{monitor.id}</span>
          <span>{monitor.type}</span>
          <span>{monitor.description}</span>
        </Flex>
      </CardContent>
    </Card>
  );
}
