import type { PushMonitor, PushMonitorConfig } from "@/features/monitors/types";
import { ConfigTable } from "@/features/monitors/components/DetailsPage/ConfigTable";
import { type FieldConfig } from "@/features/monitors/components/DetailsPage/config-fields";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts";
import { GlobalParametersApi } from "@/features/instance/global-parameters-api";
import { Flex } from "@/components/common/Flex";
import { CopyToClipboardButton } from "@/components/common/CopyToClipboardButton";
import { Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export interface PushConfigContentProps {
  monitor: PushMonitor;
}

const pushFields: FieldConfig<PushMonitorConfig> = {
  gracePeriodSeconds: {
    label: "Grace Period (s)",
  },
};

export const PushConfigContent = ({ monitor }: PushConfigContentProps) => {
  const { data: globalParams } = useQuery({
    queryKey: [QUERY_KEYS.GLOBAL_PARAMETERS],
    queryFn: () => GlobalParametersApi.getAll(),
  });

  if (!monitor.probeConfig) return null;

  const publicUrl = globalParams?.find(
    (p) => p.key === "instance.public_url",
  )?.value;
  const baseUrl = String(publicUrl || window.location.origin).replace(
    /\/+$/,
    "",
  );
  const pushUrl = `${baseUrl}/api/v1/push/${monitor.id}`;

  return (
    <Flex direction="column" className="gap-2">
      <Flex className="items-center gap-1">
        <span className="font-medium">Push URL</span>
        <Tooltip>
          <TooltipTrigger
            render={
              <Button
                variant="ghost"
                size="icon-xs"
                aria-label="Push URL usage"
              />
            }
          >
            <Info />
          </TooltipTrigger>
          <TooltipContent>
            <p>
              Add <code>?status=down</code> to report a failure and{" "}
              <code>latency=&lt;ms&gt;</code> to record latency. A POST body is
              saved as the message.
            </p>
          </TooltipContent>
        </Tooltip>
      </Flex>
      <Flex className="items-center justify-between gap-2 rounded-lg border bg-muted/50 pl-3">
        <code className="text-xs break-all">{pushUrl}</code>
        <CopyToClipboardButton value={pushUrl} size="lg" />
      </Flex>
      <ConfigTable data={monitor.probeConfig} fields={pushFields} />
    </Flex>
  );
};
