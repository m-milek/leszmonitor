import { MonitorsApi } from "@/features/monitors/monitors-api";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { Monitor } from "@/features/monitors/types";
import { TypographyH3 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { StyledLink } from "@/components/common/StyledLink";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { LucideEdit, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { StatusDot } from "@/components/common/StatusDot";
import { QUERY_KEYS } from "@/lib/consts";
import { useQuery } from "@tanstack/react-query";
import { MonitorStatusPill } from "@/features/monitors/components/MonitorStatusPill";

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
    queryFn: () => MonitorsApi.results.getLatest(monitor.id),
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
