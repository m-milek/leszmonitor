import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useSetAtom } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import { useEffect } from "react";
import { getTeam } from "@/lib/data/teamData.ts";

export const Route = createFileRoute("/_authenticated/$teamId")({
  loader: async ({ params }) => {
    const team = await getTeam(params.teamId);
    return { team };
  },
  component: TeamLayout,
});

function TeamLayout() {
  const { team } = Route.useLoaderData();
  const setTeamAtom = useSetAtom(teamAtom);

  useEffect(() => {
    if (team) {
      setTeamAtom(team);
    }
  }, [team, setTeamAtom]);

  return <Outlet />;
}
