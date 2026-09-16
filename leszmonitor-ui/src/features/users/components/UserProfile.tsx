import type { User } from "@/features/users/types";
import { TypographyH1, TypographyH2 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Initial } from "@/features/users/components/Initial";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { CopyToClipboardButton } from "@/components/common/CopyToClipboardButton";

export interface UserProfileProps {
  user: User;
}

export const UserProfile = ({ user }: UserProfileProps) => {
  return (
    <section>
      <Flex direction="column" className="gap-4">
        <Flex className="gap-4">
          <Initial text={user.username} size="xl" />
          <div className="flex flex-col justify-center">
            <TypographyH1>{user.username}</TypographyH1>
            <span className="text-muted-foreground">{user.id}</span>
          </div>
        </Flex>
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
