import { useState } from "react";
import { FileOutputIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ExportYamlDialog } from "@/components/common/ExportYamlDialog";
import { NotificationChannelsApi } from "@/features/notifications/notification-channels-api";
import type { NotificationChannel } from "@/features/notifications/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface ExportNotificationChannelButtonProps {
  channel: NotificationChannel;
}

export function ExportNotificationChannelButton({
  channel,
}: Readonly<ExportNotificationChannelButtonProps>) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <Button
        variant="ghost"
        size="icon"
        title={`Export channel ${channel.name}`}
        aria-label={`Export channel ${channel.name}`}
        onClick={() => setIsOpen(true)}
      >
        <FileOutputIcon />
      </Button>
      <ExportYamlDialog
        title={`Export ${channel.name}`}
        queryKey={[QUERY_KEYS.NOTIFICATION_CHANNELS, channel.id, "export"]}
        queryFn={() => NotificationChannelsApi.exportYaml(channel.id)}
        open={isOpen}
        onOpenChange={setIsOpen}
      />
    </>
  );
}
