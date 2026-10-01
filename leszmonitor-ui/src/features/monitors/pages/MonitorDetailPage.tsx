import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MonitorResultsList } from "@/features/monitors/components/MonitorResultsList";
import { MonitorStatusPill } from "@/features/monitors/components/MonitorStatusPill";
import { LineChart } from "@/features/monitors/components/charts/LineChartLazy";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { formatTime } from "@/features/monitors/components/charts/utils";
import type { MonitorResult } from "@/features/monitors/types";
import type { Pagination } from "@/lib/types";
import { QUERY_KEYS } from "@/lib/consts";
import { formatDuration } from "@/lib/utils.ts";
import { MonitorStatsCard } from "@/features/monitors/components/DetailsPage/MonitorStatsCard.tsx";
import { HttpStatusCodeChart } from "@/features/monitors/components/DetailsPage/HttpStatusCodeChart.tsx";
import { Center } from "@/components/common/Center.tsx";
import { TagsApi } from "@/features/tags/tags-api.ts";
import { Tag } from "@/features/tags/components/Tag.tsx";
import { MonitorConfigCard } from "@/features/monitors/components/DetailsPage/MonitorConfigCard.tsx";
import MonitorActionsGroup from "@/features/monitors/components/DetailsPage/MonitorActionsGroup.tsx";
import { useNavigate } from "@tanstack/react-router";

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

  const { data: tags } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
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
  });

  const getTagById = (id: string) => tags?.find((tag) => tag.id === id);

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

  const handleEditMonitor = () => {
    navigate({
      to: "/monitors/$monitorSlug/edit",
      params: { monitorSlug },
    });
  };

  return (
    <PageContainer>
      <Flex direction="row" className="justify-between">
        <TypographyH1>{monitor.name}</TypographyH1>
        <MonitorActionsGroup
          monitor={monitor}
          handleToggleMonitorState={handleToggleMonitorState}
          handleEditMonitor={handleEditMonitor}
          handleManuallyRunMonitor={handleManuallyRunMonitor}
          isPaused={isPaused}
        />
      </Flex>

      <Flex direction="row" className="gap-2">
        {monitor.tagIds.map((id) => {
          const tag = getTagById(id);
          return tag && <Tag tag={tag} key={id} size="lg" />;
        })}
      </Flex>

      <Flex direction="row" className="gap-4">
        <MonitorStatusPill monitor={monitor} />
      </Flex>

      <Card>
        <CardContent>
          <BatteryChart monitorResults={monitorResults ?? []} />
        </CardContent>
      </Card>

      <MonitorConfigCard monitor={monitor} />

      <Card className="min-w-0">
        <CardContent className="min-w-0">
          <Flex direction="column" className="gap-2">
            <pre className="overflow-x-auto text-xs pb-4">
              {JSON.stringify(monitor, null, 2)}
            </pre>
          </Flex>
        </CardContent>
      </Card>
      {stats && <MonitorStatsCard stats={stats} />}
      <Card>
        <CardHeader>
          <CardTitle>Statistics</CardTitle>
        </CardHeader>
        <CardContent className="min-w-0">
          {stats && (
            <>
              <p>
                {monitorStatus.toUpperCase()} for{"  "}
                {formatDuration(stats.statusChange.secondsInCurrentStatus)}
              </p>
              {stats.probeSpecific && (
                <pre className="overflow-x-auto text-xs font-mono">
                  {JSON.stringify(stats.probeSpecific, null, 2)}
                </pre>
              )}
            </>
          )}
        </CardContent>
      </Card>
      <Flex direction="row" className="gap-4 h-96 min-h-0 min-w-0 w-full">
        <Card className="flex-1 flex flex-col min-h-0 min-w-0">
          <CardHeader>
            <CardTitle>Latency</CardTitle>
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
        {monitor.type === "http" && (
          <Card>
            <CardHeader>
              <CardTitle>HTTP Status Codes</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0">
              {stats?.probeType === "http" && (
                <Center>
                  <HttpStatusCodeChart
                    data={stats.probeSpecific!.httpCodeToCount}
                  />
                </Center>
              )}
            </CardContent>
          </Card>
        )}
      </Flex>
      <Card className="flex-1 flex flex-col min-h-0 min-w-0">
        <CardContent className="flex-1 min-h-0">
          <MonitorResultsList monitor={monitor} pagination={pagination} />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
