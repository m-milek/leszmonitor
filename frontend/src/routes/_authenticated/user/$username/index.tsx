import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useQuery } from "@tanstack/react-query";
import { fetchUser } from "@/lib/data/userData.ts";
import { UserProfilePage } from "@/components/leszmonitor/UserProfilePage.tsx";

export const Route = createFileRoute("/_authenticated/user/$username/")({
  component: UserProfileComponent,
});

function UserProfileComponent() {
  const { username } = Route.useParams();

  const { data: user } = useQuery({
    queryKey: ["users", username],
    queryFn: () => fetchUser(username),
  });

  if (!user) {
    return null;
  }

  return (
    <MainPanelContainer>
      <UserProfilePage user={user} />
    </MainPanelContainer>
  );
}
