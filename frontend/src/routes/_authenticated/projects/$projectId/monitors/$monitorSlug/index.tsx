import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { getMonitorBySlug } from "@/lib/data/monitorData.ts";
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

  const { projectId, monitorSlug } = Route.useParams();

  const { data: monitor } = useQuery({
    queryKey: ["monitors", monitorSlug, projectId],
    queryFn: () => getMonitorBySlug(projectId, monitorSlug),
  });

  const largePagination: Pagination = {
    page: 1,
    perPage: 100,
  };
  const { data: monitorResults } = useQuery({
    enabled: !!monitor,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor?.id ?? "", largePagination],
    queryFn: () => getMonitorResultsByMonitorId(monitor!.id, largePagination),
  });

  if (!monitor) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>{monitor.name}</TypographyH1>
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
