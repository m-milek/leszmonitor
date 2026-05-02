import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { deleteMonitor, getMonitors } from "@/lib/data/monitorData.ts";
import { Button } from "@/components/ui/button.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { MonitorListItem } from "@/components/leszmonitor/MonitorListItem.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/",
)({
  component: MonitorsComponent,
});

function MonitorsComponent() {
  const { projectId } = Route.useParams();

  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["monitors", projectId],
    queryFn: () => getMonitors(projectId),
  });

  const { mutateAsync: deleteMutation } = useMutation({
    mutationFn: async (monitorId: string) => {
      await deleteMonitor(projectId, monitorId);
    },
  });

  const navigate = Route.useNavigate();

  if (!data) {
    return null;
  }

  const onDeleteMonitor = async (monitorId: string) => {
    await deleteMutation(monitorId);
    queryClient.invalidateQueries({ queryKey: ["monitors", projectId] });
  };

  const navigateToEditMonitor = (monitorId: string) => {
    navigate({
      to: "/projects/$projectId/monitors/$monitorSlug/edit",
      params: { projectId, monitorSlug: monitorId },
    });
  };

  return (
    <MainPanelContainer>
      <TypographyH1>Monitors</TypographyH1>
      <Card>
        <CardHeader>
          <StyledLink
            to="/projects/$projectId/monitors/new"
            params={{ projectId }}
          >
            <Button>New Monitor</Button>
          </StyledLink>
        </CardHeader>
        <CardContent>
          <Flex direction="column" className="gap-4">
            {data.map((monitor) => (
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
        <CardFooter>Footer</CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
