import type { Monitor, MonitorStatus } from "@/lib/types.ts";
import { TypographyH3 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { LucideEdit, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { StatusDot } from "@/components/leszmonitor/ui/StatusDot.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getLatestMonitorResultByMonitorId } from "@/lib/data/monitorResultsData.ts";
import { useQuery } from "@tanstack/react-query";
import { MonitorStatusPill } from "@/components/leszmonitor/MonitorStatusPill.tsx";

const monitorStatusToStatusDot = (status: MonitorStatus | undefined) => {
  switch (status) {
    case "up":
      return "success";
    case "down":
      return "failure";
    default:
      return "pending";
  }
};

export interface MonitorListItemProps {
  monitor: Monitor;
  projectSlug: string;
  onDeleteMonitor?: (monitorId: string) => Promise<void>;
  navigateToEditMonitor?: (monitorId: string) => void;
}

export function MonitorListItem({
  monitor,
  projectSlug,
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
                to="/projects/$projectId/monitors/$monitorSlug"
                params={{ projectId: projectSlug, monitorSlug: monitor.slug }}
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
