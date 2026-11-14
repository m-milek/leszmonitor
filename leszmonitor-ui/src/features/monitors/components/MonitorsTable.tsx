import { useMemo, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { type ColumnDef } from "@tanstack/react-table";
import {
  EllipsisIcon,
  LucideCirclePlay,
  LucideEdit,
  LucideLock,
  LucideTrash2,
  PauseIcon,
  PlayIcon,
} from "lucide-react";
import { DataTable, dataTableFeatures } from "@/components/common/DataTable";
import { StyledLink } from "@/components/common/StyledLink";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
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
import { QUERY_KEYS } from "@/lib/consts";
import { formatDuration } from "@/lib/utils";
import { MonitorsApi } from "@/features/monitors/monitors-api";
import type {
  MonitorStatus,
  MonitorWithStatus,
} from "@/features/monitors/types";
import { MonitorStatusBadge } from "@/features/monitors/components/MonitorStatusBadge";
import { MonitorStatePill } from "@/features/monitors/components/MonitorStatePill";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { Tag } from "@/features/tags/components/Tag";
import type { Tag as TagType } from "@/features/tags/types";

const STATUS_SORT_ORDER: MonitorStatus[] = [
  "down",
  "unknown",
  "maintenance",
  "paused",
  "up",
];

const getEffectiveStatus = (monitor: MonitorWithStatus): MonitorStatus =>
  monitor.runState === "paused" ? "paused" : monitor.status;

const StatusCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  if (monitor.runState === "paused") {
    return <MonitorStatePill monitor={monitor} />;
  }

  return (
    <MonitorStatusBadge status={monitor.status}>
      {monitor.status}
    </MonitorStatusBadge>
  );
};

const NameCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  return (
    <div className="flex items-center gap-1">
      <StyledLink
        to="/monitors/$monitorSlug"
        params={{ monitorSlug: monitor.slug }}
      >
        {monitor.name}
      </StyledLink>
      {monitor.source === "config" && (
        <Tooltip>
          <TooltipTrigger render={<span className="text-muted-foreground" />}>
            <LucideLock className="size-4" />
          </TooltipTrigger>
          <TooltipContent>Monitor defined in config file</TooltipContent>
        </Tooltip>
      )}
    </div>
  );
};

const RECENT_RESULTS_COUNT = 10;

const RecentResultsCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  const { data: results = [] } = useQuery({
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id, "recent"],
    queryFn: () =>
      MonitorsApi.results.getPage(monitor.id, {
        page: 1,
        perPage: RECENT_RESULTS_COUNT,
      }),
  });

  return (
    <div className="w-40">
      <BatteryChart
        monitorResults={results ?? []}
        config={{ barWidth: 16, heightClassName: "h-6", showTimeRange: false }}
      />
    </div>
  );
};

const ActionsCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const isPaused = monitor.runState === "paused";
  const isConfigBased = monitor.source === "config";

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

  return (
    <div className="flex items-center justify-end">
      {monitor.type !== "push" && (
        <Tooltip>
          <TooltipTrigger
            render={
              <Button
                variant="ghost"
                size="icon"
                aria-label="Run now"
                onClick={() => runMutation.mutate()}
                disabled={runMutation.isPending}
              />
            }
          >
            <LucideCirclePlay />
          </TooltipTrigger>
          <TooltipContent>Run now</TooltipContent>
        </Tooltip>
      )}
      <Tooltip>
        <TooltipTrigger
          render={
            <Button
              variant="ghost"
              size="icon"
              aria-label={isPaused ? "Resume" : "Pause"}
              onClick={() => toggleStateMutation.mutate()}
              disabled={isConfigBased || toggleStateMutation.isPending}
            />
          }
        >
          {isPaused ? <PlayIcon /> : <PauseIcon />}
        </TooltipTrigger>
        <TooltipContent>{isPaused ? "Resume" : "Pause"}</TooltipContent>
      </Tooltip>
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <Button variant="ghost" size="icon" aria-label="More actions" />
          }
        >
          <EllipsisIcon />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItemIcon
            icon={LucideEdit}
            onClick={() =>
              navigate({
                to: "/monitors/$monitorSlug/edit",
                params: { monitorSlug: monitor.slug },
              })
            }
            disabled={isConfigBased}
          >
            Edit
          </DropdownMenuItemIcon>
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
      <DeleteMonitorDialog
        monitor={monitor}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  );
};

export interface MonitorsTableProps {
  monitors: MonitorWithStatus[];
  tags: TagType[];
}

export const MonitorsTable = ({ monitors, tags }: MonitorsTableProps) => {
  const columns = useMemo<
    ColumnDef<typeof dataTableFeatures, MonitorWithStatus>[]
  >(
    () => [
      {
        accessorKey: "name",
        header: "Name",
        sortFn: (a, b) =>
          a.original.name.localeCompare(b.original.name, undefined, {
            sensitivity: "base",
          }),
        cell: ({ row }) => <NameCell monitor={row.original} />,
      },
      {
        id: "status",
        header: "Status",
        accessorFn: getEffectiveStatus,
        sortFn: (a, b) =>
          STATUS_SORT_ORDER.indexOf(getEffectiveStatus(a.original)) -
          STATUS_SORT_ORDER.indexOf(getEffectiveStatus(b.original)),
        cell: ({ row }) => <StatusCell monitor={row.original} />,
      },
      {
        id: "recentResults",
        header: "Recent checks",
        cell: ({ row }) => <RecentResultsCell monitor={row.original} />,
      },
      {
        id: "tags",
        header: "Tags",
        cell: ({ row }) => (
          <div className="flex flex-wrap gap-1">
            {tags
              .filter((tag) => row.original.tagIds?.includes(tag.id))
              .map((tag) => (
                <Tag key={tag.id} tag={tag} />
              ))}
          </div>
        ),
      },
      {
        accessorKey: "type",
        header: "Type",
        sortFn: (a, b) => a.original.type.localeCompare(b.original.type),
        cell: ({ row }) => (
          <Badge variant="outline">{row.original.type.toUpperCase()}</Badge>
        ),
      },
      {
        accessorKey: "interval",
        header: "Interval",
        sortFn: (a, b) => a.original.interval - b.original.interval,
        cell: ({ row }) => (
          <span className="text-sm text-muted-foreground">
            {formatDuration(row.original.interval)}
          </span>
        ),
      },
      {
        id: "actions",
        cell: ({ row }) => <ActionsCell monitor={row.original} />,
      },
    ],
    [tags],
  );

  return (
    <DataTable
      data={monitors}
      columns={columns}
      compact
      getRowId={(monitor) => monitor.id}
      sortable
      emptyMessage="No monitors."
    />
  );
};
