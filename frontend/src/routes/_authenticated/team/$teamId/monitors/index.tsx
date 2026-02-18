import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/Typography.tsx";
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

export const Route = createFileRoute("/_authenticated/team/$teamId/monitors/")({
  component: MonitorsComponent,
});

function MonitorsComponent() {
  const { teamId } = Route.useParams();

  const { data } = useQuery({
    queryKey: ["monitors", teamId],
    queryFn: () => getMonitors(teamId),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Monitors</TypographyH1>
      <Card>
        <CardHeader>
          <StyledLink to="/team/$teamId/monitors/new" params={{ teamId }}>
            <Button>New Monitor</Button>
          </StyledLink>
        </CardHeader>
        <CardContent>{JSON.stringify(data)}</CardContent>
        <CardFooter>Footer</CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
