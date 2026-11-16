import { useQuery } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { CopyToClipboardButton } from "@/components/common/CopyToClipboardButton";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import type { StatusPage } from "@/features/status-pages/types";
import { QUERY_KEYS } from "@/lib/consts";

export interface ExportStatusPageDialogProps {
  statusPage: StatusPage;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const ExportStatusPageContent = ({
  statusPage,
}: {
  statusPage: StatusPage;
}) => {
  const { data: yaml, isError } = useQuery({
    queryKey: [QUERY_KEYS.STATUS_PAGES, statusPage.id, "export"],
    queryFn: () => StatusPagesApi.exportYaml(statusPage.id),
    staleTime: 0,
    gcTime: 0,
  });

  if (isError) {
    return (
      <span className="text-destructive">
        Failed to export the status page.
      </span>
    );
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

export const ExportStatusPageDialog = ({
  statusPage,
  open,
  onOpenChange,
}: ExportStatusPageDialogProps) => (
  <Dialog open={open} onOpenChange={onOpenChange}>
    <DialogContent className="sm:max-w-2xl">
      <DialogHeader>
        <DialogTitle>Export {statusPage.title}</DialogTitle>
      </DialogHeader>
      <ExportStatusPageContent statusPage={statusPage} />
    </DialogContent>
  </Dialog>
);
