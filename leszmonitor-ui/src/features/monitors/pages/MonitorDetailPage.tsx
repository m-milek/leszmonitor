import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  LucideCirclePlay,
  PauseIcon,
  PencilIcon,
  PlayIcon,
} from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1, TypographyH2 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { ButtonGroup } from "@/components/ui/button-group";
import { Button } from "@/components/ui/button";
import { MonitorResultsList } from "@/features/monitors/components/MonitorResultsList";
import { MonitorStatusPill } from "@/features/monitors/components/MonitorStatusPill";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog";
import { LineChart } from "@/features/monitors/components/charts/LineChartLazy";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { formatTime } from "@/features/monitors/components/charts/utils";
import type { MonitorResult, MonitorStatus } from "@/features/monitors/types";
import type { Pagination } from "@/lib/types";
import { QUERY_KEYS } from "@/lib/consts";
import { formatDuration } from "@/lib/utils.ts";

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
    queryFn: () => MonitorsApi.getBySlug(monitorSlug),
  });

  const { data: monitorResults } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", pagination],
    queryFn: () => MonitorsApi.results.getPage(monitor!.id, pagination),
  });

  const { data: stats } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_LATENCY_STATS, monitor?.id ?? ""],
    queryFn: () =>
      MonitorsApi.stats.get(monitor!.id, {
        from: new Date(Date.now() - 24 * 60 * 60 * 1000), // last 24 hours
      }),
  });

  const mutation = useMutation({
    mutationKey: [QUERY_KEYS.MONITORS, monitorSlug],
    mutationFn: async () =>
      MonitorsApi.updateState(monitor!.id, isPaused ? "active" : "paused"),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
      });
    },
  });

  const manuallyRunMutation = useMutation({
    mutationKey: [QUERY_KEYS.MONITORS, monitorSlug, "run"],
    mutationFn: async () => MonitorsApi.run(monitor!.id),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", pagination],
      });
    },
  });

  if (!monitor) {
    return null;
  }

  const isPaused = monitor.runState === "paused";

  const monitorStatus = monitorResults?.[0]?.status ?? "unknown";

  const handleToggleMonitorState = () => {
    mutation.mutate();
  };

  const handleManuallyRunMonitor = () => {
    manuallyRunMutation.mutate();
  };

  const statusCounts = stats?.uptime.statusToCount ?? {};
  const statusPercentages = Object.entries(
    stats?.uptime.statusToPercentage ?? {},
  );

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
            size="icon-lg"
            onClick={handleToggleMonitorState}
          >
            {isPaused ? <PlayIcon /> : <PauseIcon />}
          </Button>
          <Button variant="outline" size="icon-lg" onClick={handleEditMonitor}>
            <PencilIcon />
          </Button>
          <Button
            variant="outline"
            size="icon-lg"
            onClick={handleManuallyRunMonitor}
          >
            <LucideCirclePlay />
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
          <TypographyH2>Statistics</TypographyH2>
          {stats && (
            <>
              <p>Avg: {stats.latency.avg.toFixed(2)} ms</p>
              <p>Min: {stats.latency.min.toFixed(2)} ms</p>
              <p>Max: {stats.latency.max.toFixed(2)} ms</p>
              <p>
                {monitorStatus.toUpperCase()} for{"  "}
                {formatDuration(stats.statusChange.secondsInCurrentStatus)}
              </p>
              {statusPercentages.map(([status, percentage]) => (
                <p key={status}>
                  {status.toUpperCase()}: {percentage.toFixed(2)}% (
                  {statusCounts[status as MonitorStatus] ?? 0})
                </p>
              ))}
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
