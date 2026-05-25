import { createFileRoute } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getMonitorBySlug,
  updateMonitorState,
} from "@/lib/data/monitorData.ts";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import type { Pagination } from "@/lib/types.ts";
import { MonitorResultsList } from "@/components/leszmonitor/MonitorResultsList.tsx";
import { LatencyChart } from "@/components/leszmonitor/charts/LatencyChart.tsx";
import { getMonitorResultsByMonitorId } from "@/lib/data/monitorResultsData.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { ButtonGroup } from "@/components/ui/button-group.tsx";
import { Button } from "@/components/ui/button.tsx";
import { PauseIcon, PlayIcon, TrashIcon } from "lucide-react";
import { MonitorStatusPill } from "@/components/leszmonitor/MonitorStatusPill.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/$monitorSlug/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const pagination: Pagination = {
    page: 1,
    perPage: 20,
  };

  const largePagination: Pagination = {
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
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", largePagination],
    queryFn: () => getMonitorResultsByMonitorId(monitor!.id, largePagination),
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

  const isPaused = monitor.state === "paused";

  const handleToggleMonitorState = () => {
    mutation.mutate();
  };

  return (
    <MainPanelContainer>
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

      <Card>
        <CardContent>
          <pre>{JSON.stringify(monitor, null, 2)}</pre>
        </CardContent>
      </Card>
      <Flex direction="row" className="gap-4 h-96 min-h-0">
        <Card className="flex-1 flex flex-col min-h-0">
          <CardHeader>
            <TypographyH2>Latency (ms)</TypographyH2>
          </CardHeader>
          <CardContent className="flex-1 min-h-0">
            <LatencyChart monitorResults={monitorResults ?? []} />
          </CardContent>
        </Card>

        <Card className="flex-1 flex flex-col min-h-0">
          <CardContent className="flex-1 min-h-0">
            <MonitorResultsList monitor={monitor} pagination={pagination} />
          </CardContent>
        </Card>
      </Flex>
    </MainPanelContainer>
  );
}
