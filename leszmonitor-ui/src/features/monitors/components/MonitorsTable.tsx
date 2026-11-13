import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useMutation, useQueryClient } from "@tanstack/react-query";
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
import { StatusDot } from "@/components/common/StatusDot";
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
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { MonitorWithStatus } from "@/features/monitors/types";
import { MonitorStatusBadge } from "@/features/monitors/components/MonitorStatusBadge";
import { MonitorStatePill } from "@/features/monitors/components/MonitorStatePill";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog";
import { Tag } from "@/features/tags/components/Tag";
import type { Tag as TagType } from "@/features/tags/types";

const StatusCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  if (monitor.runState === "paused") {
    return (
      <div className="flex items-center gap-2">
        <StatusDot status="paused" />
        <MonitorStatePill monitor={monitor} />
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <StatusDot status={monitorStatusToStatusDot(monitor.status)} />
      <MonitorStatusBadge status={monitor.status}>
        {monitor.status}
      </MonitorStatusBadge>
    </div>
  );
};

const NameCell = ({ monitor }: { monitor: MonitorWithStatus }) => {
  const link = (
    <StyledLink
      to="/monitors/$monitorSlug"
      params={{ monitorSlug: monitor.slug }}
    >
      {monitor.name}
    </StyledLink>
  );

  return (
    <div className="flex items-center gap-1">
      {monitor.description ? (
        <Tooltip>
          <TooltipTrigger render={<span />}>{link}</TooltipTrigger>
          <TooltipContent>{monitor.description}</TooltipContent>
        </Tooltip>
      ) : (
        link
      )}
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
  const columns: ColumnDef<typeof dataTableFeatures, MonitorWithStatus>[] = [
    {
      id: "status",
      header: "Status",
      cell: ({ row }) => <StatusCell monitor={row.original} />,
    },
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => <NameCell monitor={row.original} />,
    },
    {
      accessorKey: "type",
      header: "Type",
      cell: ({ row }) => (
        <Badge variant="outline">{row.original.type.toUpperCase()}</Badge>
      ),
    },
    {
      accessorKey: "interval",
      header: "Interval",
      cell: ({ row }) => (
        <span className="text-sm text-muted-foreground">
          {formatDuration(row.original.interval)}
        </span>
      ),
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
      id: "actions",
      cell: ({ row }) => <ActionsCell monitor={row.original} />,
    },
  ];

  return (
    <DataTable
      data={monitors}
      columns={columns}
      compact
      getRowId={(monitor) => monitor.id}
      emptyMessage="No monitors."
    />
  );
};
