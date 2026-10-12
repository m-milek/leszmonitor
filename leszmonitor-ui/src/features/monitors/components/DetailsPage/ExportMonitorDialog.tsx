import { useQuery } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { CopyToClipboardButton } from "@/components/common/CopyToClipboardButton";
import { MonitorsApi } from "@/features/monitors/monitors-api";
import { QUERY_KEYS } from "@/lib/consts";
import type { Monitor } from "@/features/monitors/types";

export interface ExportMonitorDialogProps {
  monitor: Monitor;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const ExportMonitorContent = ({ monitor }: { monitor: Monitor }) => {
  const { data: yaml, isError } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitor.id, "export"],
    queryFn: () => MonitorsApi.exportYaml(monitor.id),
    staleTime: 0,
    gcTime: 0,
  });

  if (isError) {
    return (
      <span className="text-destructive">Failed to export the monitor.</span>
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

export const ExportMonitorDialog = ({
  monitor,
  open,
  onOpenChange,
}: ExportMonitorDialogProps) => (
  <Dialog open={open} onOpenChange={onOpenChange}>
    <DialogContent className="sm:max-w-2xl">
      <DialogHeader>
        <DialogTitle>Export {monitor.name}</DialogTitle>
      </DialogHeader>
      <ExportMonitorContent monitor={monitor} />
    </DialogContent>
  </Dialog>
);
