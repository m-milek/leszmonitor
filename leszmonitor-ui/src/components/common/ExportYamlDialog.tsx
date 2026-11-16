import { type QueryKey, useQuery } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { CopyToClipboardButton } from "@/components/common/CopyToClipboardButton";

export interface ExportYamlDialogProps {
  title: string;
  queryKey: QueryKey;
  queryFn: () => Promise<string>;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const ExportYamlContent = ({
  queryKey,
  queryFn,
}: Pick<ExportYamlDialogProps, "queryKey" | "queryFn">) => {
  const { data: yaml, isError } = useQuery({
    queryKey,
    queryFn,
    staleTime: 0,
    gcTime: 0,
  });

  if (isError) {
    return <span className="text-destructive">Export failed.</span>;
  }

  if (yaml === undefined) {
    return <Skeleton className="h-64 w-full" />;
  }

  return (
    <div className="relative min-w-0">
      <pre className="max-h-[70vh] overflow-auto rounded-lg bg-muted p-4 font-mono text-xs">
        {yaml}
      </pre>
      <div className="absolute top-2 right-2">
        <CopyToClipboardButton value={yaml} />
      </div>
    </div>
  );
};

export const ExportYamlDialog = ({
  title,
  queryKey,
  queryFn,
  open,
  onOpenChange,
}: ExportYamlDialogProps) => (
  <Dialog open={open} onOpenChange={onOpenChange}>
    <DialogContent className="sm:max-w-2xl">
      <DialogHeader>
        <DialogTitle>{title}</DialogTitle>
      </DialogHeader>
      <ExportYamlContent queryKey={queryKey} queryFn={queryFn} />
    </DialogContent>
  </Dialog>
);
