import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "@/components/ui/toast";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  StatusPagesApi,
  type StatusPagePayload,
} from "@/features/status-pages/status-pages-api";
import { StatusPageForm } from "@/features/status-pages/forms/StatusPageForm";
import type { StatusPage } from "@/features/status-pages/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface EditStatusPageDialogProps {
  statusPage: StatusPage;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function EditStatusPageDialog({
  statusPage,
  open,
  onOpenChange,
}: EditStatusPageDialogProps) {
  const queryClient = useQueryClient();
  const formId = `edit-status-page-form-${statusPage.id}`;

  const updateMutation = useMutation({
    mutationFn: (values: StatusPagePayload) =>
      StatusPagesApi.update(statusPage.id, values),
    onSuccess: (page) => {
      toast.add({
        title: `Status page "${page.title}" updated`,
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
          <DialogTitle>Edit Status Page</DialogTitle>
        </DialogHeader>
        <StatusPageForm
          key={String(open)}
          id={formId}
          defaultValues={statusPage}
          onSubmit={async (values) => {
            await updateMutation.mutateAsync(values);
          }}
        />
        <DialogFooter>
          <DialogClose render={<Button variant="outline" />}>Close</DialogClose>
          <Button
            type="submit"
            form={formId}
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? "Saving..." : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
