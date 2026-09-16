import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type {
  Monitor,
  MonitorResult,
} from "@/features/monitors/types";
import type { Pagination } from "@/lib/types";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts";
import { getMonitorResultsByMonitorId } from "@/features/monitors/results-api";
import { Flex } from "@/components/common/Flex";
import { StatusDot } from "@/components/common/StatusDot";
import { ScrollArea } from "@/components/ui/scroll-area";

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
