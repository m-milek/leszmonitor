import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Trash2 } from "lucide-react";
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
import { Tag } from "@/components/leszmonitor/Tag.tsx";
import { deleteTag } from "@/lib/data/tags-api.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import type { Tag as TagModel } from "@/lib/types.ts";

export interface DeleteTagDialogProps {
  tag: TagModel;
}

export function DeleteTagDialog({ tag }: Readonly<DeleteTagDialogProps>) {
  const queryClient = useQueryClient();
  const [isOpen, setIsOpen] = useState(false);

  const deleteTagMutation = useMutation({
    mutationFn: () => deleteTag(tag.id),
    onSuccess: () => {
      toast.success(`Tag "${tag.name}" deleted`);
      setIsOpen(false);
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.TAGS] });
    },
    onError: (error) => {
      toast.error("Failed to delete tag: " + error.message);
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          title={`Delete tag ${tag.name}`}
          aria-label={`Delete tag ${tag.name}`}
        >
          <Trash2 className="text-destructive" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete tag</DialogTitle>
          <DialogDescription>
            This removes the tag from every monitor using it. This cannot be
            undone.
          </DialogDescription>
        </DialogHeader>
        <div className="flex justify-center py-2">
          <Tag tag={tag} />
        </div>
        <DialogFooter showCloseButton>
          <Button
            variant="destructive"
            onClick={() => deleteTagMutation.mutate()}
            disabled={deleteTagMutation.isPending}
          >
            {deleteTagMutation.isPending ? "Deleting..." : "Delete Tag"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
