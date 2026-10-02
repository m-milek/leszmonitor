import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MonitorResultsCard } from "@/features/monitors/components/DetailsPage/MonitorResultsCard";
import { LineChart } from "@/features/monitors/components/charts/LineChartLazy";
import { Skeleton } from "@/components/ui/skeleton";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { formatTime } from "@/features/monitors/components/charts/utils";
import type { MonitorResult } from "@/features/monitors/types";
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
import { MonitorStatusBadge } from "@/features/monitors/components/MonitorStatusBadge";
import { MonitorStatePill } from "@/features/monitors/components/MonitorStatePill.tsx";
import {
  resultsRangePerPage,
  ResultsRangeSelect,
  resultsRangeStart,
} from "@/features/monitors/components/DetailsPage/ResultsRangeSelect";

export interface MonitorDetailPageProps {
  monitorSlug: string;
  range: number;
  onRangeChange: (range: number) => void;
}

const latencyChartConfig = {
  durationMs: {
    label: "Latency",
  },
};

export function MonitorDetailPage({
  monitorSlug,
  range,
  onRangeChange,
}: MonitorDetailPageProps) {
  const navigate = useNavigate();

  const queryClient = useQueryClient();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
    queryFn: () => MonitorsApi.getBySlug(monitorSlug),
  });

  const { data: monitorResults } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", range],
    queryFn: () =>
      MonitorsApi.results.getPage(
        monitor!.id,
        { page: 1, perPage: resultsRangePerPage(range, monitor!.interval) },
        resultsRangeStart(range),
      ),
  });

  const { data: stats } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_LATENCY_STATS, monitor?.id ?? "", range],
    queryFn: () =>
      MonitorsApi.stats.get(monitor!.id, { from: resultsRangeStart(range) }),
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
      <Flex
        direction="row"
        directionMobile="column"
        className="justify-between max-md:gap-2"
      >
        <TypographyH1>{monitor.name}</TypographyH1>
        <MonitorActionsGroup
          monitor={monitor}
          handleToggleMonitorState={handleToggleMonitorState}
          handleEditMonitor={handleEditMonitor}
          handleManuallyRunMonitor={handleManuallyRunMonitor}
          isPaused={isPaused}
        />
      </Flex>

      <Flex
        direction="row"
        directionMobile="column"
        className="justify-between max-md:gap-2"
      >
        <Flex direction="row" className="gap-2">
          <Flex direction="row" className="gap-2">
            {monitorResults && stats ? (
              <MonitorStatusBadge status={monitorResults[0]?.status} size="lg">
                {monitorStatus.toUpperCase()} for{" "}
                {formatDuration(stats.statusChange.secondsInCurrentStatus)}
              </MonitorStatusBadge>
            ) : (
              <Skeleton className="h-(--text-3xl) w-40 rounded-4xl" />
            )}
            {monitor.runState === "paused" && (
              <MonitorStatePill monitor={monitor} size="lg" />
            )}
          </Flex>
          <Flex direction="row" className="gap-2">
            {monitor.tagIds.map((id) => {
              const tag = getTagById(id);
              return tag && <Tag tag={tag} key={id} size="lg" />;
            })}
          </Flex>
        </Flex>

        <Flex
          direction="row"
          className="gap-4 items-center max-md:justify-between"
        >
          <span>Runs every {formatDuration(monitor.interval)}</span>
          <ResultsRangeSelect value={range} onChange={onRangeChange} />
        </Flex>
      </Flex>

      {monitor.description?.length !== 0 && (
        <span className="text-muted-foreground">{monitor.description}</span>
      )}

      <Card>
        <CardContent>
          {monitorResults ? (
            <BatteryChart monitorResults={monitorResults} />
          ) : (
            <Skeleton className="h-8 w-full" />
          )}
        </CardContent>
      </Card>

      <MonitorConfigCard monitor={monitor} />

      <MonitorStatsCard stats={stats} />

      <Flex
        direction="row"
        directionMobile="column"
        className="gap-4 h-96 max-md:h-auto min-h-0 min-w-0 w-full"
      >
        <Card className="flex-2 flex flex-col min-h-0 min-w-0 max-md:flex-none max-md:h-80">
          <CardHeader>
            <CardTitle>Latency</CardTitle>
          </CardHeader>
          <CardContent className="flex-1 min-h-0">
            {monitorResults ? (
              <LineChart<MonitorResult>
                data={monitorResults}
                config={latencyChartConfig}
                timestampExtractor={(r) => new Date(r.createdAt).getTime()}
                xAxisKey="createdAt"
                yAxisKey="durationMs"
                uniqueMatchKey="id"
                xAxisTickFormatter={formatTime}
                yAxisDomain={[0, "auto"]}
              />
            ) : (
              <Skeleton className="h-full w-full" />
            )}
          </CardContent>
        </Card>
        {monitor.type === "http" && (
          <Card className="flex-1 min-w-0 max-md:flex-none">
            <CardHeader>
              <CardTitle>HTTP Status Codes</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0">
              {!stats ? (
                <Skeleton className="aspect-video h-64" />
              ) : (
                stats.probeType === "http" && (
                  <Center>
                    <HttpStatusCodeChart
                      data={stats.probeSpecific!.httpCodeToCount}
                      expectedStatusCodes={
                        monitor.probeConfig?.expectedStatusCodes ?? []
                      }
                    />
                  </Center>
                )
              )}
            </CardContent>
          </Card>
        )}
      </Flex>
      <MonitorResultsCard
        key={range}
        results={monitorResults === null ? [] : monitorResults}
      />
    </PageContainer>
  );
}
