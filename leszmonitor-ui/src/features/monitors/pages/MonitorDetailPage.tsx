import { useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PauseIcon, PencilIcon, PlayIcon } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1, TypographyH2 } from "@/components/common/Typography.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { ButtonGroup } from "@/components/ui/button-group.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  getMonitorBySlug,
  updateMonitorState,
} from "@/features/monitors/api/monitors.ts";
import { getMonitorResultsByMonitorId } from "@/features/monitors/api/results.ts";
import { getLatencyStatsByMonitorId } from "@/features/monitors/api/stats.ts";
import { MonitorResultsList } from "@/features/monitors/components/MonitorResultsList.tsx";
import { MonitorStatusPill } from "@/features/monitors/components/MonitorStatusPill.tsx";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog.tsx";
import { LineChart } from "@/features/monitors/components/charts/LineChartLazy.tsx";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart.tsx";
import { formatTime } from "@/features/monitors/components/charts/utils.ts";
import type { MonitorResult } from "@/features/monitors/model/types.ts";
import type { Pagination } from "@/lib/types.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";

export interface MonitorDetailPageProps {
  monitorSlug: string;
}

const latencyChartConfig = {
  durationMs: {
    label: "Latency (ms)",
  },
};

export function MonitorDetailPage({ monitorSlug }: MonitorDetailPageProps) {
  const pagination: Pagination = {
    page: 1,
    perPage: 100,
  };

  const navigate = useNavigate();

  const queryClient = useQueryClient();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
    queryFn: () => getMonitorBySlug(monitorSlug),
  });

  const { data: monitorResults } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", pagination],
    queryFn: () => getMonitorResultsByMonitorId(monitor!.id, pagination),
  });

  const { data: latencyStats } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_LATENCY_STATS, monitor?.id ?? ""],
    queryFn: () =>
      getLatencyStatsByMonitorId(monitor!.id, {
        from: new Date(Date.now() - 24 * 60 * 60 * 1000), // last 24 hours
      }),
  });

  const mutation = useMutation({
    mutationKey: [QUERY_KEYS.MONITORS, monitorSlug],
    mutationFn: async () =>
      updateMonitorState(monitor!.id, isPaused ? "active" : "paused"),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
      });
    },
  });

  if (!monitor) {
    return null;
  }

  const isPaused = monitor.runState === "paused";

  const handleToggleMonitorState = () => {
    mutation.mutate();
  };

  const handleEditMonitor = () => {
    navigate({
      to: "/monitors/$monitorSlug/edit",
      params: { monitorSlug },
    });
  };

  return (
    <PageContainer>
      <TypographyH1>{monitor.name}</TypographyH1>
      <Flex direction="row" className="gap-4">
        <ButtonGroup>
          <Button
            variant="outline"
            className="size-10"
            onClick={handleToggleMonitorState}
          >
            {isPaused ? <PlayIcon /> : <PauseIcon />}
          </Button>
          <Button
            variant="outline"
            className="size-10"
            onClick={handleEditMonitor}
          >
            <PencilIcon />
          </Button>
          <DeleteMonitorDialog
            monitor={monitor}
            onDeleted={() => navigate({ to: "/monitors" })}
          />
        </ButtonGroup>
        <MonitorStatusPill monitor={monitor} />
      </Flex>

      <Card className="min-w-0">
        <CardContent className="min-w-0">
          <pre className="overflow-x-auto text-xs pb-4">
            {JSON.stringify(monitor, null, 2)}
          </pre>
          <BatteryChart monitorResults={monitorResults ?? []} />
        </CardContent>
      </Card>
      <Card>
        <CardContent className="min-w-0">
          <TypographyH2>Latency (last 24h)</TypographyH2>
          {latencyStats && (
            <>
              <p>Avg: {latencyStats.averageLatency.toFixed(2)} ms</p>
              <p>Min: {latencyStats.minLatency.toFixed(2)} ms</p>
              <p>Max: {latencyStats.maxLatency.toFixed(2)} ms</p>
            </>
          )}
        </CardContent>
      </Card>
      <Flex direction="row" className="gap-4 h-96 min-h-0 min-w-0 w-full">
        <Card className="flex-1 flex flex-col min-h-0 min-w-0">
          <CardHeader>
            <TypographyH2>Latency (ms)</TypographyH2>
          </CardHeader>
          <CardContent className="flex-1 min-h-0">
            <LineChart<MonitorResult>
              data={monitorResults ?? []}
              config={latencyChartConfig}
              timestampExtractor={(r) => new Date(r.createdAt).getTime()}
              xAxisKey="createdAt"
              yAxisKey="durationMs"
              uniqueMatchKey="id"
              xAxisTickFormatter={formatTime}
              yAxisDomain={[0, "auto"]}
            />
          </CardContent>
        </Card>

        <Card className="flex-1 flex flex-col min-h-0 min-w-0">
          <CardContent className="flex-1 min-h-0">
            <MonitorResultsList monitor={monitor} pagination={pagination} />
          </CardContent>
        </Card>
      </Flex>
    </PageContainer>
  );
}
