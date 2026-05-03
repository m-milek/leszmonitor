import type { Monitor } from "@/lib/types.ts";
import { TypographyH3 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { LucideEdit, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { StatusDot } from "@/components/leszmonitor/ui/StatusDot.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getLatestMonitorResultById } from "@/lib/data/monitorResultsData.ts";
import { useQuery } from "@tanstack/react-query";

export interface MonitorListItemProps {
  monitor: Monitor;
  projectId: string;
  onDeleteMonitor: (monitorId: string) => Promise<void>;
  navigateToEditMonitor: (monitorId: string) => void;
}

export function MonitorListItem({
  monitor,
  projectId,
  onDeleteMonitor,
  navigateToEditMonitor,
}: MonitorListItemProps) {
  const { data: lastResultData } = useQuery({
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id],
    queryFn: () => getLatestMonitorResultById(monitor.id),
  });

  const dotStatus = lastResultData
    ? lastResultData.isSuccess
      ? "success"
      : "failure"
    : "pending";

  return (
    <Card>
      <CardHeader>
        <Flex direction="row" className="justify-between">
          <Flex direction="row" className="items-center gap-2">
            <StatusDot status={dotStatus} />
            <TypographyH3>
              <StyledLink
                to="/projects/$projectId/monitors/$monitorSlug"
                params={{ projectId, monitorSlug: monitor.slug }}
              >
                {monitor.name}
              </StyledLink>
            </TypographyH3>
          </Flex>
          <Flex direction="row">
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
              onClick={() => onDeleteMonitor(monitor.slug)}
            >
              <LucideTrash2 className="size-5 text-destructive" />
            </Button>
          </Flex>
        </Flex>
      </CardHeader>
      <CardContent>
        <Flex direction="column">
          <span>{monitor.id}</span>
          <span>{monitor.type}</span>
          <span>{monitor.description}</span>
          {lastResultData ? (
            <span>
              Last result: {lastResultData.isSuccess ? "Success" : "Failure"} at{" "}
              {new Date(lastResultData.createdAt).toLocaleString()}
            </span>
          ) : (
            <span>No results yet</span>
          )}
        </Flex>
      </CardContent>
    </Card>
  );
}
