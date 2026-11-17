import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Pencil } from "lucide-react";
import { toast } from "@/components/ui/toast";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  NotificationChannelsApi,
  type NotificationChannelPayload,
} from "@/features/notifications/notification-channels-api";
import { NotificationChannelForm } from "@/features/notifications/forms/NotificationChannelForm";
import { channelToFormValues } from "@/features/notifications/schema";
import type { NotificationChannel } from "@/features/notifications/types";
import { QUERY_KEYS } from "@/lib/consts";

const FORM_ID = "edit-notification-channel-form";

export interface EditNotificationChannelDialogProps {
  channel: NotificationChannel;
}

export function EditNotificationChannelDialog({
  channel,
}: Readonly<EditNotificationChannelDialogProps>) {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const updateMutation = useMutation({
    mutationFn: (values: NotificationChannelPayload) =>
      NotificationChannelsApi.update(channel.id, values),
    onSuccess: (updated) => {
      toast.add({
        title: `Notification channel "${updated.name}" updated`,
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
            title={`Edit channel ${channel.name}`}
            aria-label={`Edit channel ${channel.name}`}
            disabled={channel.source === "config"}
          />
        }
      >
        <Pencil />
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Notification Channel</DialogTitle>
          <DialogDescription>
            Secrets are not shown. Leave them empty to keep the current values.
          </DialogDescription>
        </DialogHeader>
        <NotificationChannelForm
          key={String(isOpen)}
          id={FORM_ID}
          defaultValues={channelToFormValues(channel)}
          isEdit
          onSubmit={async (values) => {
            await updateMutation.mutateAsync(values);
          }}
        />
        <DialogFooter>
          <DialogClose render={<Button variant="outline" />}>Close</DialogClose>
          <Button
            type="submit"
            form={FORM_ID}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? "Saving..." : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
