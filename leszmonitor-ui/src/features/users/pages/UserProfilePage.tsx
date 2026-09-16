import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { UserProfile } from "@/features/users/components/UserProfile.tsx";
import { getUser } from "@/features/users/users-api.ts";

export interface UserProfilePageProps {
  username: string;
}

export function UserProfilePage({ username }: UserProfilePageProps) {
  const { data: user } = useQuery({
    queryKey: ["users", username],
    queryFn: () => getUser(username),
  });

  if (!user) {
    return null;
  }

  return (
    <PageContainer>
      <UserProfile user={user} />
    </PageContainer>
  );
}
