import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { MonitorWithStatus } from "@/features/monitors/types";
import { TypographyH3 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { StyledLink } from "@/components/common/StyledLink";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { LucideEdit, LucideLock, LucideTrash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { StatusDot } from "@/components/common/StatusDot";
import { QUERY_KEYS } from "@/lib/consts";
import { useQuery } from "@tanstack/react-query";
import { Tag } from "@/features/tags/components/Tag.tsx";
import { TagsApi } from "@/features/tags/tags-api.ts";
import { MonitorStatusBadge } from "@/features/monitors/components/MonitorStatusBadge.tsx";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export interface MonitorListItemProps {
  monitor: MonitorWithStatus;
  onDeleteMonitor?: (monitorId: string) => Promise<void>;
  navigateToEditMonitor?: (monitorId: string) => void;
}

export function MonitorListItem({
  monitor,
  onDeleteMonitor,
  navigateToEditMonitor,
}: Readonly<MonitorListItemProps>) {
  const { data: tags } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
  });

  if (!tags) {
    return null;
  }

  const dotStatus = monitorStatusToStatusDot(monitor.status);
  const isConfigBased = monitor.source === "config";

  return (
    <Card className="transition-colors hover:bg-muted/50">
      <CardHeader>
        <Flex direction="row" className="justify-between">
          <Flex direction="row" className="items-center gap-2">
            <StatusDot status={dotStatus} />
            <TypographyH3>
              <StyledLink
                to="/monitors/$monitorSlug"
                params={{ monitorSlug: monitor.slug }}
              >
                {monitor.name}
              </StyledLink>
            </TypographyH3>
            <MonitorStatusBadge status={monitor.status}>
              {monitor.status}
            </MonitorStatusBadge>
            {isConfigBased && (
              <Tooltip>
                <TooltipTrigger>
                  <Badge variant="ghost" className="h-7 px-2 [&>svg]:size-5!">
                    <LucideLock />
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>Monitor defined in config file</TooltipContent>
              </Tooltip>
            )}
          </Flex>
          <Flex direction="row">
            {navigateToEditMonitor && onDeleteMonitor && (
              <>
                <Button
                  variant="ghost"
                  size="icon-lg"
                  onClick={() => navigateToEditMonitor(monitor.slug)}
                  disabled={isConfigBased}
                >
                  <LucideEdit className="size-5" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon-lg"
                  onClick={() => onDeleteMonitor(monitor.id)}
                  disabled={isConfigBased}
                >
                  <LucideTrash2 className="size-5 text-destructive" />
                </Button>
              </>
            )}
          </Flex>
        </Flex>
      </CardHeader>
      <CardContent>
        <Flex direction="column" className="gap-2">
          <Flex direction="row" className="gap-2">
            {tags
              .filter((tag) => monitor.tagIds?.includes(tag.id))
              .map((tag) => (
                <Tag key={tag.id} tag={tag} />
              ))}
          </Flex>
          <span>{monitor.description}</span>
        </Flex>
      </CardContent>
    </Card>
  );
}
