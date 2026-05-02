import type { Monitor } from "@/lib/types.ts";
import { TypographyH3 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { LucideEdit, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";

export interface MonitorListItemProps {
  monitor: Monitor;
  projectId: string;
  onDeleteMonitor: (monitorId: string) => Promise<void>;
  navigateToEditMonitor: (monitorId: string) => void;
}

export function MonitorListItem({
  monitor,
  projectId,
  onDeleteMonitor,
  navigateToEditMonitor,
}: MonitorListItemProps) {
  return (
    <Card>
      <CardHeader>
        <Flex direction="row" className="justify-between">
          <TypographyH3>
            <StyledLink
              to="/projects/$projectId/monitors/$monitorSlug"
              params={{ projectId, monitorSlug: monitor.slug }}
            >
              {monitor.name}
            </StyledLink>
          </TypographyH3>
          <div>
            <Button
              variant="ghost"
              size="icon-lg"
              onClick={() => navigateToEditMonitor(monitor.slug)}
            >
              <LucideEdit className="size-6" />
            </Button>
            <Button
              variant="ghost"
              size="icon-lg"
              onClick={() => onDeleteMonitor(monitor.slug)}
            >
              <LucideTrash2 className="size-6 text-destructive" />
            </Button>
          </div>
        </Flex>
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
