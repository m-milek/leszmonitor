import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useAtomValue } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/sidebar/Typography.tsx";
import { TeamMembersTable } from "@/components/leszmonitor/tables/TeamMembersTable.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId/members/")({
  component: TeamRoute,
});

function TeamRoute() {
  const team = useAtomValue(teamAtom);

  if (!team) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Members</TypographyH1>
      <Card>
        <CardHeader>
          <TypographyH2>
            {team.members.length}{" "}
            {team.members.length === 1 ? "Member" : "Members"}
          </TypographyH2>
        </CardHeader>
        <CardContent>
          <TeamMembersTable members={team.members} />
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
