import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useSetAtom } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import { useEffect } from "react";
import { getTeam } from "@/lib/data/teamData.ts";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_authenticated/team/$teamId")({
  component: TeamLayout,
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
