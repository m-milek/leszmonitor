import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { TrashIcon } from "lucide-react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { Button } from "@/components/ui/button.tsx";
import { deleteMonitor } from "@/features/monitors/monitors-api.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import type { Monitor } from "@/features/monitors/types.ts";

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
    mutationFn: () => deleteMonitor(monitor.id),
    onSuccess: () => {
      toast.success(`Monitor "${monitor.name}" deleted`);
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
      onDeleted?.();
    },
    onError: (error) => {
      toast.error("Failed to delete monitor: " + error.message);
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button
          variant="destructive"
          className="size-10"
          title={`Delete monitor ${monitor.name}`}
          aria-label={`Delete monitor ${monitor.name}`}
        >
          <TrashIcon />
        </Button>
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
