import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useSetAtom } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import { useEffect } from "react";
import { getTeam } from "@/lib/data/teamData.ts";
import { useQuery } from "@tanstack/react-query";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId")({
  component: TeamLayout,
  notFoundComponent: NotFound,
});

function TeamLayout() {
  const { teamId } = Route.useParams();
  const setTeamAtom = useSetAtom(teamAtom);

  const { data: team } = useQuery({
    queryKey: ["team", teamId],
    queryFn: () => getTeam(teamId),
  });

  useEffect(() => {
    if (team) {
      setTeamAtom(team);
    }
  }, [team, teamId, setTeamAtom]);

  return <Outlet />;
}

function NotFound() {
  return (
    <MainPanelContainer>
      <TypographyH1>Not Found</TypographyH1>;
    </MainPanelContainer>
  );
}
