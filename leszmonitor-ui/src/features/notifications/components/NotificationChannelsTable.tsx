import { type ColumnDef } from "@tanstack/react-table";
import { LucideLock } from "lucide-react";
import { DataTable, dataTableFeatures } from "@/components/common/DataTable";
import { Flex } from "@/components/common/Flex";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { formatDate } from "@/lib/utils";
import {
  notificationChannelKindLabels,
  type NotificationChannel,
} from "@/features/notifications/types";
import { EditNotificationChannelDialog } from "@/features/notifications/components/EditNotificationChannelDialog";
import { DeleteNotificationChannelDialog } from "@/features/notifications/components/DeleteNotificationChannelDialog";
import { ExportNotificationChannelButton } from "@/features/notifications/components/ExportNotificationChannelButton";

const channelTarget = (channel: NotificationChannel): string => {
  switch (channel.kind) {
    case "webhook":
      return channel.config.url;
    case "ntfy":
      return `${channel.config.serverUrl.replace(/\/$/, "")}/${channel.config.topic}`;
    case "gotify":
      return channel.config.serverUrl;
  }
};

export interface NotificationChannelsTableProps {
  channels: NotificationChannel[];
}

const columns: ColumnDef<typeof dataTableFeatures, NotificationChannel>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => (
      <Flex direction="row" className="items-center gap-2">
        <span className="font-medium">{row.original.name}</span>
        {row.original.source === "config" && (
          <Tooltip>
            <TooltipTrigger>
              <Badge variant="ghost" className="h-7 px-2 [&>svg]:size-5!">
                <LucideLock />
              </Badge>
            </TooltipTrigger>
            <TooltipContent>Channel defined in config file</TooltipContent>
          </Tooltip>
        )}
      </Flex>
    ),
  },
  {
    accessorKey: "kind",
    header: "Type",
    cell: ({ row }) => (
      <Badge variant="outline">
        {notificationChannelKindLabels[row.original.kind]}
      </Badge>
    ),
  },
  {
    id: "target",
    header: "Target",
    cell: ({ row }) => (
      <code className="block max-w-80 truncate text-muted-foreground">
        {channelTarget(row.original)}
      </code>
    ),
  },
  {
    accessorKey: "enabled",
    header: "Status",
    cell: ({ row }) =>
      row.original.enabled ? (
        <Badge variant="secondary">Enabled</Badge>
      ) : (
        <Badge variant="outline" className="text-muted-foreground">
          Disabled
        </Badge>
      ),
  },
  {
    accessorKey: "createdAt",
    header: "Created",
    cell: ({ row }) => (
      <span className="whitespace-nowrap">
        {formatDate(row.original.createdAt)}
      </span>
    ),
  },
  {
    id: "actions",
    header: "",
    cell: ({ row }) => (
      <div className="flex justify-end">
        <ExportNotificationChannelButton channel={row.original} />
        <EditNotificationChannelDialog channel={row.original} />
        <DeleteNotificationChannelDialog channel={row.original} />
      </div>
    ),
  },
];

export const NotificationChannelsTable = ({
  channels,
}: NotificationChannelsTableProps) => {
  return (
    <DataTable
      data={channels}
      columns={columns}
      emptyMessage="No notification channels yet."
    />
  );
};
