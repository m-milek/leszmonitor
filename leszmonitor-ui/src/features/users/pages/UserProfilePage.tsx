import { UsersApi } from "@/features/users/users-api";
import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { UserProfile } from "@/features/users/components/UserProfile";

export interface UserProfilePageProps {
  username: string;
}

export function UserProfilePage({ username }: UserProfilePageProps) {
  const { data: user } = useQuery({
    queryKey: ["users", username],
    queryFn: () => UsersApi.get(username),
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
