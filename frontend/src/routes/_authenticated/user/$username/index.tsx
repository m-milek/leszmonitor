import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/Typography.tsx";
import { useQuery } from "@tanstack/react-query";
import { fetchUser } from "@/lib/data/userData.ts";
import { Card, CardContent } from "@/components/ui/card.tsx";

export const Route = createFileRoute("/_authenticated/user/$username/")({
  component: UserProfileComponent,
});

function UserProfileComponent() {
  const { username } = Route.useParams();

  const { data } = useQuery({
    queryKey: ["users", username],
    queryFn: () => fetchUser(username),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>{username}</TypographyH1>
      <Card>
        <CardContent>{JSON.stringify(data, null, 2)}</CardContent>
      </Card>
    </MainPanelContainer>
  );
}
