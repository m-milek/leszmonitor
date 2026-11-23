import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Trash2 } from "lucide-react";
import { toast } from "@/components/ui/toast";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { NotificationChannelsApi } from "@/features/notifications/notification-channels-api";
import type { NotificationChannel } from "@/features/notifications/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface DeleteNotificationChannelDialogProps {
  channel: NotificationChannel;
}

export function DeleteNotificationChannelDialog({
  channel,
}: Readonly<DeleteNotificationChannelDialogProps>) {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const deleteMutation = useMutation({
    mutationFn: () => NotificationChannelsApi.remove(channel.id),
    onSuccess: () => {
      toast.add({
        title: `Notification channel "${channel.name}" deleted`,
        type: "success",
      });
      setIsOpen(false);
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.NOTIFICATION_CHANNELS],
      });
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          <Button
            variant="ghost"
            size="icon"
            title={`Delete channel ${channel.name}`}
            aria-label={`Delete channel ${channel.name}`}
            disabled={channel.source === "config"}
          />
        }
      >
        <Trash2 className="text-destructive" />
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete notification channel</DialogTitle>
          <DialogDescription>
            Monitors using &quot;{channel.name}&quot; will stop sending
            notifications to it. This cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter showCloseButton>
          <Button
            variant="destructive"
            onClick={() => deleteMutation.mutate()}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? "Deleting..." : "Delete Channel"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
