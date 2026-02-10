import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useAtomValue } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import { TypographyH1 } from "@/components/leszmonitor/sidebar/Typography.tsx";

export const Route = createFileRoute("/_authenticated/$teamId/team/")({
  component: TeamRoute,
});

function TeamRoute() {
  const team = useAtomValue(teamAtom);

  if (!team) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>{team.name}</TypographyH1>
      {team.members.length === 0 ? (
        <p className="mt-2 text-gray-500">No members in this team yet.</p>
      ) : (
        <ul>
          {team.members.map((member) => (
            <li key={member.id}>
              <div className="flex items-center gap-2">
                <div>{member.id}</div>
                <div>{member.role}</div>
                <div>{member.createdAt.toISOString()}</div>
                <div>{member.updatedAt.toISOString()}</div>
              </div>
            </li>
          ))}
        </ul>
      )}
    </MainPanelContainer>
  );
}
