import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MonitorResultsCard } from "@/features/monitors/components/DetailsPage/MonitorResultsCard";
import { LineChart } from "@/features/monitors/components/charts/LineChartLazy";
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
  type ResultsRange,
  ResultsRangeSelect,
  resultsRangePerPage,
  resultsRangeStart,
} from "@/features/monitors/components/DetailsPage/ResultsRangeSelect";
import { useState } from "react";

export interface MonitorDetailPageProps {
  monitorSlug: string;
}

const latencyChartConfig = {
  durationMs: {
    label: "Latency",
  },
};

export function MonitorDetailPage({ monitorSlug }: MonitorDetailPageProps) {
  const [range, setRange] = useState<ResultsRange>("24h");

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
        <Flex direction="row" className="gap-4 items-center">
          <TypographyH1>{monitor.name}</TypographyH1>
        </Flex>
        <MonitorActionsGroup
          monitor={monitor}
          handleToggleMonitorState={handleToggleMonitorState}
          handleEditMonitor={handleEditMonitor}
          handleManuallyRunMonitor={handleManuallyRunMonitor}
          isPaused={isPaused}
        />
      </Flex>

      <Flex direction="row" className="justify-between">
        <Flex direction="row" className="gap-2">
          <Flex direction="row" className="gap-2">
            <MonitorStatusBadge status={monitorResults?.[0]?.status} size="lg">
              {monitorStatus.toUpperCase()} for{" "}
              {formatDuration(stats?.statusChange.secondsInCurrentStatus ?? 0)}
            </MonitorStatusBadge>
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

        <Flex direction="row" className="gap-4 items-center">
          <span>Runs every {formatDuration(monitor.interval)}</span>
          <ResultsRangeSelect value={range} onChange={setRange} />
        </Flex>
      </Flex>

      {monitor.description?.length !== 0 && (
        <span className="text-muted-foreground">{monitor.description}</span>
      )}

      <Card>
        <CardContent>
          <BatteryChart monitorResults={monitorResults ?? []} />
        </CardContent>
      </Card>

      <MonitorConfigCard monitor={monitor} />

      {stats && <MonitorStatsCard stats={stats} />}

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
      <MonitorResultsCard results={monitorResults ?? []} />
    </PageContainer>
  );
}
