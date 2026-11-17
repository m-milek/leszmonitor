import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { Plus } from "lucide-react";
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

const FORM_ID = "new-notification-channel-form";

export function NewNotificationChannelDialog() {
  const [isOpen, setIsOpen] = useState(false);

  const createMutation = useMutation({
    mutationFn: (values: NotificationChannelPayload) =>
      NotificationChannelsApi.create(values),
    onSuccess: (channel) => {
      toast.add({
        title: `Notification channel "${channel.name}" created`,
        type: "success",
      });
      setIsOpen(false);
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger render={<Button size="sm" />}>
        <Plus />
        Add Channel
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Notification Channel</DialogTitle>
          <DialogDescription>
            Where to send alerts when a monitor changes status.
          </DialogDescription>
        </DialogHeader>
        <NotificationChannelForm
          key={String(isOpen)}
          id={FORM_ID}
          onSubmit={async (values) => {
            await createMutation.mutateAsync(values);
          }}
        />
        <DialogFooter>
          <DialogClose render={<Button variant="outline" />}>Close</DialogClose>
          <Button
            type="submit"
            form={FORM_ID}
            disabled={createMutation.isPending}
          >
            {createMutation.isPending ? "Creating..." : "Create Channel"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
