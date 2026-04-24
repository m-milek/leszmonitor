import type { Monitor } from "@/lib/types.ts";
import { TypographyH3 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";

export interface MonitorListItemProps {
  monitor: Monitor;
  projectId: string;
}

export function MonitorListItem({ monitor, projectId }: MonitorListItemProps) {
  return (
    <Card>
      <CardHeader>
        <TypographyH3>
          <StyledLink
            to="/projects/$projectId/monitors/$monitorSlug"
            params={{ projectId, monitorSlug: monitor.slug }}
          >
            {monitor.name}
          </StyledLink>
        </TypographyH3>
      </CardHeader>
      <CardContent>
        <Flex direction="column">
          <span>{monitor.id}</span>
          <span>{monitor.type}</span>
          <span>{monitor.description}</span>
        </Flex>
      </CardContent>
    </Card>
  );
}
