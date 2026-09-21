import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { TrashIcon } from "lucide-react";
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
import { QUERY_KEYS } from "@/lib/consts";
import type { Monitor } from "@/features/monitors/types";

export interface DeleteMonitorDialogProps {
  monitor: Monitor;
  onDeleted?: () => void;
}

export function DeleteMonitorDialog({
  monitor,
  onDeleted,
}: Readonly<DeleteMonitorDialogProps>) {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const deleteMonitorMutation = useMutation({
    mutationFn: () => MonitorsApi.remove(monitor.id),
    onSuccess: () => {
      toast.add({
        title: `Monitor "${monitor.name}" deleted`,
        type: "success",
      });
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
      onDeleted?.();
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger
        render={
          <Button
            variant="destructive"
            size="icon-lg"
            title={`Delete monitor ${monitor.name}`}
            aria-label={`Delete monitor ${monitor.name}`}
          />
        }
      >
        <TrashIcon />
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete monitor</DialogTitle>
          <DialogDescription>
            This permanently removes &quot;{monitor.name}&quot; and all of its
            results. This cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter showCloseButton>
          <Button
            variant="destructive"
            onClick={() => deleteMonitorMutation.mutate()}
            disabled={deleteMonitorMutation.isPending}
          >
            {deleteMonitorMutation.isPending ? "Deleting..." : "Delete Monitor"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
