import { monitorStatusToStatusDot } from "@/features/monitors/model/status.ts";
import type {
  Monitor,
  MonitorResult,
} from "@/features/monitors/model/types.ts";
import type { Pagination } from "@/lib/types.ts";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getMonitorResultsByMonitorId } from "@/features/monitors/api/results.ts";
import { Flex } from "@/components/common/Flex.tsx";
import { StatusDot } from "@/components/common/StatusDot.tsx";
import { ScrollArea } from "@/components/ui/scroll-area.tsx";

export interface MonitorResultsListProps {
  monitor: Monitor;
  pagination: Pagination;
}

export const MonitorResultsList = ({
  monitor,
  pagination,
}: MonitorResultsListProps) => {
  const { data: results } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id, pagination],
    queryFn: () => getMonitorResultsByMonitorId(monitor.id, pagination),
  });

  const getStatusDotStatus = (result: MonitorResult) => {
    return monitorStatusToStatusDot(result.status);
  };

  const getStatusText = (result: MonitorResult) => {
    return result.status.toUpperCase();
  };

  return (
    <div className="h-full">
      <ScrollArea className="h-full">
        {results?.map((result) => (
          <div key={result.id}>
            <Flex direction="row" className="gap-4 items-center">
              <StatusDot status={getStatusDotStatus(result)} />
              <span className="font-mono">{result.id.substring(0, 8)}</span>
              <span className="font-mono">{getStatusText(result)}</span>
              <span>{result.createdAt.toLocaleString()}</span>
            </Flex>
          </div>
        ))}
      </ScrollArea>
    </div>
  );
};
