import { useState } from "react";
import { cn } from "cn";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { Monitor, MonitorWithStatus } from "@/features/monitors/types";
import { TypographyH3 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { StyledLink } from "@/components/common/StyledLink";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  EllipsisIcon,
  LucideCirclePlay,
  LucideEdit,
  LucideLock,
  LucideTrash2,
  PauseIcon,
  PlayIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { StatusDot } from "@/components/common/StatusDot";
import { QUERY_KEYS } from "@/lib/consts";
import { formatDuration } from "@/lib/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Tag } from "@/features/tags/components/Tag";
import { TagsApi } from "@/features/tags/tags-api";
import { MonitorsApi } from "@/features/monitors/monitors-api";
import { MonitorStatusBadge } from "@/features/monitors/components/MonitorStatusBadge";
import { MonitorStatePill } from "@/features/monitors/components/MonitorStatePill";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export interface MonitorListItemProps {
  monitor: MonitorWithStatus;
  navigateToEditMonitor?: (monitorSlug: string) => void;
}

const getMonitorTarget = (monitor: Monitor): string | undefined => {
  switch (monitor.type) {
    case "http":
      return monitor.probeConfig?.url;
    case "tcp":
      return (
        monitor.probeConfig &&
        `${monitor.probeConfig.host}:${monitor.probeConfig.port}`
      );
    case "dns":
      return (
        monitor.probeConfig &&
        `${monitor.probeConfig.recordType} ${monitor.probeConfig.hostname}`
      );
    case "ping":
      return monitor.probeConfig?.host;
    case "push":
      return undefined;
  }
};

export function MonitorListItem({
  monitor,
  navigateToEditMonitor,
}: Readonly<MonitorListItemProps>) {
  const queryClient = useQueryClient();
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  const { data: tags } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
  });

  const isPaused = monitor.runState === "paused";

  const toggleStateMutation = useMutation({
    mutationFn: () =>
      MonitorsApi.updateState(monitor.id, isPaused ? "active" : "paused"),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
    },
  });

  const runMutation = useMutation({
    mutationFn: () => MonitorsApi.run(monitor.id),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id, "latest"],
      });
    },
  });

  if (!tags) {
    return null;
  }

  const dotStatus = isPaused
    ? "paused"
    : monitorStatusToStatusDot(monitor.status);
  const isConfigBased = monitor.source === "config";
  const target = getMonitorTarget(monitor);

  return (
    <Card
      className={cn(
        "transition-colors hover:bg-muted/50",
        isPaused && "opacity-60",
      )}
    >
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
            {isPaused ? (
              <MonitorStatePill monitor={monitor} />
            ) : (
              <MonitorStatusBadge status={monitor.status}>
                {monitor.status}
              </MonitorStatusBadge>
            )}
            {isConfigBased && (
              <Tooltip>
                <TooltipTrigger>
                  <Badge variant="ghost" className="h-7 px-2 [&>svg]:size-5!">
                    <LucideLock />
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>Monitor defined in config file</TooltipContent>
              </Tooltip>
            )}
          </Flex>
          <Flex direction="row" className="items-center">
            {monitor.type !== "push" && (
              <Tooltip>
                <TooltipTrigger
                  render={
                    <Button
                      variant="ghost"
                      size="icon-lg"
                      aria-label="Run now"
                      onClick={() => runMutation.mutate()}
                      disabled={runMutation.isPending}
                    />
                  }
                >
                  <LucideCirclePlay className="size-5" />
                </TooltipTrigger>
                <TooltipContent>Run now</TooltipContent>
              </Tooltip>
            )}
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon-lg"
                    aria-label={isPaused ? "Resume" : "Pause"}
                    onClick={() => toggleStateMutation.mutate()}
                    disabled={isConfigBased || toggleStateMutation.isPending}
                  />
                }
              >
                {isPaused ? (
                  <PlayIcon className="size-5" />
                ) : (
                  <PauseIcon className="size-5" />
                )}
              </TooltipTrigger>
              <TooltipContent>{isPaused ? "Resume" : "Pause"}</TooltipContent>
            </Tooltip>
            <DropdownMenu>
              <DropdownMenuTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon-lg"
                    aria-label="More actions"
                  />
                }
              >
                <EllipsisIcon className="size-5" />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {navigateToEditMonitor && (
                  <DropdownMenuItemIcon
                    icon={LucideEdit}
                    onClick={() => navigateToEditMonitor(monitor.slug)}
                    disabled={isConfigBased}
                  >
                    Edit
                  </DropdownMenuItemIcon>
                )}
                <DropdownMenuItemIcon
                  icon={LucideTrash2}
                  variant="destructive"
                  onClick={() => setIsDeleteOpen(true)}
                  disabled={isConfigBased}
                >
                  Delete
                </DropdownMenuItemIcon>
              </DropdownMenuContent>
            </DropdownMenu>
          </Flex>
        </Flex>
      </CardHeader>
      <CardContent>
        <Flex direction="column" className="gap-2">
          <Flex
            direction="row"
            className="items-center gap-2 text-sm text-muted-foreground"
          >
            <Badge variant="outline">{monitor.type.toUpperCase()}</Badge>
            {target && <span className="truncate font-mono">{target}</span>}
            <span>· every {formatDuration(monitor.interval)}</span>
          </Flex>
          <Flex direction="row" className="gap-2">
            {tags
              .filter((tag) => monitor.tagIds?.includes(tag.id))
              .map((tag) => (
                <Tag key={tag.id} tag={tag} />
              ))}
          </Flex>
          <span>{monitor.description}</span>
        </Flex>
      </CardContent>
      <DeleteMonitorDialog
        monitor={monitor}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </Card>
  );
}
