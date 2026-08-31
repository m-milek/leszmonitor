import { createFileRoute } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getMonitorBySlug,
  updateMonitorState,
} from "@/lib/data/monitorData.ts";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import type { MonitorResult, Pagination } from "@/lib/types.ts";
import { MonitorResultsList } from "@/components/leszmonitor/MonitorResultsList.tsx";
import { LineChart } from "@/components/leszmonitor/charts/LineChart.tsx";
import { BatteryChart } from "@/components/leszmonitor/charts/BatteryChart/BatteryChart.tsx";
import { formatTime } from "@/components/leszmonitor/charts/utils.ts";
import { getMonitorResultsByMonitorId } from "@/lib/data/monitorResultsData.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { ButtonGroup } from "@/components/ui/button-group.tsx";
import { Button } from "@/components/ui/button.tsx";
import { PauseIcon, PlayIcon, TrashIcon } from "lucide-react";
import { MonitorStatusPill } from "@/components/leszmonitor/MonitorStatusPill.tsx";
import { getLatencyStatsByMonitorId } from "@/lib/data/monitorStats.ts";

const latencyChartConfig = {
  durationMs: {
    label: "Latency (ms)",
  },
};

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/$monitorSlug/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const pagination: Pagination = {
    page: 1,
    perPage: 100,
  };

  const { projectId, monitorSlug } = Route.useParams();

  const queryClient = useQueryClient();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug, projectId],
    queryFn: () => getMonitorBySlug(projectId, monitorSlug),
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
    mutationKey: [QUERY_KEYS.MONITORS, monitorSlug, projectId],
    mutationFn: async () =>
      updateMonitorState(monitor!.id, isPaused ? "active" : "paused"),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.MONITORS, monitorSlug, projectId],
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
          <Button variant="destructive" className="size-10">
            <TrashIcon />
          </Button>
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
