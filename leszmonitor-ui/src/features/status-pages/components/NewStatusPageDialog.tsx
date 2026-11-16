import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
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
  StatusPagesApi,
  type StatusPagePayload,
} from "@/features/status-pages/status-pages-api";
import { StatusPageForm } from "@/features/status-pages/forms/StatusPageForm";
import { QUERY_KEYS } from "@/lib/consts";

const FORM_ID = "new-status-page-form";

export function NewStatusPageDialog() {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const createMutation = useMutation({
    mutationFn: (values: StatusPagePayload) => StatusPagesApi.create(values),
    onSuccess: (page) => {
      toast.add({
        title: `Status page "${page.title}" created`,
        type: "success",
      });
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STATUS_PAGES] });
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger render={<Button size="sm" />}>
        <Plus />
        New Status Page
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New Status Page</DialogTitle>
          <DialogDescription>
            A status page shows every monitor with any of the selected tags.
          </DialogDescription>
        </DialogHeader>
        <StatusPageForm
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
            {createMutation.isPending ? "Creating..." : "Create Status Page"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
