import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "@/components/ui/toast";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import type { StatusPage } from "@/features/status-pages/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface DeleteStatusPageDialogProps {
  statusPage: StatusPage;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function DeleteStatusPageDialog({
  statusPage,
  open,
  onOpenChange,
}: DeleteStatusPageDialogProps) {
  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    mutationFn: () => StatusPagesApi.remove(statusPage.id),
    onSuccess: () => {
      toast.add({
        title: `Status page "${statusPage.title}" deleted`,
        type: "success",
      });
      onOpenChange(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STATUS_PAGES] });
    },
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete status page</DialogTitle>
          <DialogDescription>
            Delete {`"${statusPage.title}"`}? This cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter showCloseButton>
          <Button
            variant="destructive"
            onClick={() => deleteMutation.mutate()}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? "Deleting..." : "Delete Status Page"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
