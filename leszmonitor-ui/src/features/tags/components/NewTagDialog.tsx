import { TagsApi } from "@/features/tags/tags-api";
import { type TagPayload } from "@/features/tags/tags-api";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { toast } from "sonner";
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
import { TagForm } from "@/features/tags/forms/TagForm";
import { QUERY_KEYS } from "@/lib/consts";

const FORM_ID = "new-tag-form";

export function NewTagDialog() {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const createTagMutation = useMutation({
    mutationFn: (values: TagPayload) => TagsApi.create(values),
    onSuccess: (tag) => {
      toast.success(`Tag "${tag.name}" created`);
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.TAGS] });
    },
    onError: (error) => {
      toast.error("Failed to create tag: " + error.message);
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Tag
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Tag</DialogTitle>
          <DialogDescription>
            Tags let you group monitors across this Leszmonitor instance.
          </DialogDescription>
        </DialogHeader>
        <TagForm
          key={String(isOpen)}
          id={FORM_ID}
          onSubmit={async (values) => {
            await createTagMutation.mutateAsync(values);
          }}
        />
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Close</Button>
          </DialogClose>
          <Button
            type="submit"
            form={FORM_ID}
            disabled={createTagMutation.isPending}
          >
            {createTagMutation.isPending ? "Creating..." : "Create Tag"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
