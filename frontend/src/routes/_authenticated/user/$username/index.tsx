import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { useQuery } from "@tanstack/react-query";
import { getUser } from "@/lib/data/userData.ts";
import { UserProfilePage } from "@/components/leszmonitor/UserProfilePage.tsx";

export const Route = createFileRoute("/_authenticated/user/$username/")({
  component: UserProfileComponent,
});

function UserProfileComponent() {
  const { username } = Route.useParams();

  const { data: user } = useQuery({
    queryKey: ["users", username],
    queryFn: () => getUser(username),
  });

  if (!user) {
    return null;
  }

  return (
    <PageContainer>
      <UserProfilePage user={user} />
    </PageContainer>
  );
}
