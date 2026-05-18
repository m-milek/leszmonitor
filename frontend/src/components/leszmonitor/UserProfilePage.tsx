import type { User } from "@/lib/types.ts";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { UserInitial } from "@/components/leszmonitor/UserInitial.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { CopyToClipboardButton } from "@/components/leszmonitor/CopyToClipboardButton.tsx";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { getProjects } from "@/lib/data/projectData.ts";
import { Suspense } from "react";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { ProjectsTable } from "@/components/leszmonitor/tables/ProjectsTable.tsx";

export interface UserProfilePageProps {
  user: User;
}

export const UserProfilePage = ({ user }: UserProfilePageProps) => {
  console.log("user", user);
  const { data: userProjects, isLoading } = useQuery({
    queryKey: [QUERY_KEYS.PROJECTS, user.username],
    enabled: !!user,
    queryFn: () => getProjects(user.username),
  });

  return (
    <section>
      <Flex direction="column" className="gap-4">
        <Flex className="gap-4">
          <UserInitial username={user.username} size="xl" />
          <div className="flex flex-col justify-center">
            <TypographyH1>{user.username}</TypographyH1>
            <span className="text-muted-foreground">{user.id}</span>
          </div>
        </Flex>
        <Card>
          <CardHeader>
            <TypographyH2>Projects</TypographyH2>
          </CardHeader>
          <CardContent>
            <Suspense fallback={<Skeleton className="h-24" />}>
              {userProjects && <ProjectsTable projects={userProjects} />}
              {isLoading && <Skeleton className="h-24" />}
            </Suspense>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <TypographyH2>Details</TypographyH2>
          </CardHeader>
          <CardContent>
            <dl className="flex flex-col gap-2">
              <div className="flex items-center gap-4">
                <dt className="text-muted-foreground w-32">ID</dt>
                <dd className="font-mono flex-1 text-right">{user.id}</dd>
                <CopyToClipboardButton value={user.id} />
              </div>

              <div className="flex items-center gap-4">
                <dt className="text-muted-foreground w-32">Username</dt>
                <dd className="font-mono flex-1 text-right">{user.username}</dd>
                <CopyToClipboardButton value={user.username} />
              </div>

              <div className="flex items-center gap-4">
                <dt className="text-muted-foreground w-32">Joined</dt>
                <dd className="font-mono flex-1 text-right">
                  {new Date(user.createdAt).toLocaleString()}
                </dd>
                <div className="w-8" />
              </div>

              <div className="flex items-center gap-4">
                <dt className="text-muted-foreground w-32">Last login</dt>
                <dd className="font-mono flex-1 text-right">TODO</dd>
                <div className="w-8" />
              </div>
            </dl>
          </CardContent>
        </Card>
      </Flex>
    </section>
  );
};
