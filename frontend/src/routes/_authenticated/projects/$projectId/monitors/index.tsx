import { createFileRoute, Link } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  deleteMonitor,
  getMonitorsByProjectSlug,
} from "@/lib/data/monitorData.ts";
import { MonitorListItem } from "@/components/leszmonitor/MonitorListItem.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { Button } from "@/components/ui/button.tsx";
import { LucidePlusCircle } from "lucide-react";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/",
)({
  component: MonitorsComponent,
});

function MonitorsComponent() {
  const { projectId } = Route.useParams();
  const queryClient = useQueryClient();

  const { data: monitors = [] } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, projectId],
    queryFn: () => getMonitorsByProjectSlug(projectId),
  });

  const { mutateAsync: deleteMutation } = useMutation({
    mutationFn: (monitorId: string) => deleteMonitor(projectId, monitorId),
  });

  const navigate = Route.useNavigate();

  const onDeleteMonitor = async (monitorId: string) => {
    await deleteMutation(monitorId);
    queryClient.invalidateQueries({
      queryKey: [QUERY_KEYS.MONITORS, projectId],
    });
  };

  const navigateToEditMonitor = (monitorSlug: string) => {
    navigate({
      to: "/projects/$projectId/monitors/$monitorSlug/edit",
      params: { projectId, monitorSlug },
    });
  };

  return (
    <MainPanelContainer>
      <TypographyH1>Monitors</TypographyH1>
      <Card>
        <CardHeader>
          <Link to={"/projects/$projectId/monitors/new"} params={{ projectId }}>
            <Button>
              <LucidePlusCircle />
              <span>New Monitor</span>
            </Button>
          </Link>
        </CardHeader>
        <CardContent>
          <Flex direction="column" className="gap-4">
            {monitors.map((monitor) => (
              <MonitorListItem
                key={monitor.id}
                monitor={monitor}
                projectId={projectId}
                onDeleteMonitor={onDeleteMonitor}
                navigateToEditMonitor={navigateToEditMonitor}
              />
            ))}
          </Flex>
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
