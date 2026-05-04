import type { Monitor, Pagination } from "@/lib/types.ts";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getMonitorResultsByMonitorId } from "@/lib/data/monitorResultsData.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { StatusDot } from "@/components/leszmonitor/ui/StatusDot.tsx";

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
  return (
    <div>
      {results?.map((result) => (
        <div key={result.id}>
          <Flex direction="row" className="gap-4">
            <StatusDot status={result.isSuccess ? "success" : "failure"} />
            <span className="font-mono">{result.id}</span>
            <span className="font-mono">
              {result.isSuccess ? "Success" : "Failure"}
            </span>
            <span>{result.createdAt.toLocaleString()}</span>
          </Flex>
        </div>
      ))}
    </div>
  );
};
