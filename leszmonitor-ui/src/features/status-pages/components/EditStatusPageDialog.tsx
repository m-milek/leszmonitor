import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { LucidePencil } from "lucide-react";
import { toast } from "@/components/ui/toast";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  StatusPagesApi,
  type StatusPagePayload,
} from "@/features/status-pages/status-pages-api";
import { StatusPageForm } from "@/features/status-pages/forms/StatusPageForm";
import type { StatusPage } from "@/features/status-pages/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface EditStatusPageDialogProps {
  statusPage: StatusPage;
}

export function EditStatusPageDialog({
  statusPage,
}: EditStatusPageDialogProps) {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);
  const formId = `edit-status-page-form-${statusPage.id}`;

  const updateMutation = useMutation({
    mutationFn: (values: StatusPagePayload) =>
      StatusPagesApi.update(statusPage.id, values),
    onSuccess: (page) => {
      toast.add({
        title: `Status page "${page.title}" updated`,
        type: "success",
      });
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STATUS_PAGES] });
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <Tooltip>
        <TooltipTrigger
          render={
            <DialogTrigger
              render={
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Edit status page"
                />
              }
            />
          }
        >
          <LucidePencil />
        </TooltipTrigger>
        <TooltipContent>Edit status page</TooltipContent>
      </Tooltip>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Status Page</DialogTitle>
        </DialogHeader>
        <StatusPageForm
          key={String(isOpen)}
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
