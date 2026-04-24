import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { useQuery } from "@tanstack/react-query";
import { getMonitors } from "@/lib/data/monitorData.ts";
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

  const { data } = useQuery({
    queryKey: ["monitors", projectId],
    queryFn: () => getMonitors(projectId),
  });

  if (!data) {
    return null;
  }

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
                monitor={monitor}
                projectId={projectId}
                key={monitor.id}
              />
            ))}
          </Flex>
        </CardContent>
        <CardFooter>Footer</CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
