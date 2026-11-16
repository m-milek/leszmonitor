import type { MetadataResponse } from "@/features/instance/metadata-api";
import { LucideInfo } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface MetadataProps {
  data?: MetadataResponse;
}

const MetadataTable = ({ data }: { data?: MetadataResponse }) => {
  const rows = [
    ["Version", data?.version],
    ["CI Build", data?.ciBuildNumber],
    ["Git Commit", data?.gitCommit],
    ["Image Tag", data?.imageTag],
  ];

  return (
    <dl className="grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm">
      {rows.map(([label, value]) => (
        <div key={label} className="contents">
          <dt className="text-muted-foreground">{label}</dt>
          <dd className="font-mono text-xs leading-5 break-all">
            {value || "—"}
          </dd>
        </div>
      ))}
    </dl>
  );
};

export const Metadata = ({ data }: MetadataProps) => {
  return (
    <Popover>
      <PopoverTrigger
        aria-label="Build information"
        className="group/metadata flex w-full cursor-pointer items-center justify-between rounded-md px-2 py-1 text-xs text-sidebar-foreground/60 transition-colors outline-hidden hover:bg-sidebar-accent hover:text-sidebar-foreground focus-visible:ring-2 focus-visible:ring-sidebar-ring data-popup-open:bg-sidebar-accent data-popup-open:text-sidebar-foreground"
      >
        <span className="truncate">
          Leszmonitor {data && (data.version || "dev")}
        </span>
        <LucideInfo className="size-3.5 shrink-0 opacity-0 transition-opacity group-hover/metadata:opacity-100 group-focus-visible/metadata:opacity-100 group-data-popup-open/metadata:opacity-100" />
      </PopoverTrigger>
      <PopoverContent side="right" align="end" sideOffset={12} className="w-80">
        <MetadataTable data={data} />
      </PopoverContent>
    </Popover>
  );
};
