import type { Monitor, MonitorType } from "@/features/monitors/types";
import {
  Card,
  CardAction,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { formatDate } from "@/lib/utils";
import { HttpConfigContent } from "@/features/monitors/components/DetailsPage/HttpConfigContent";
import { TcpConfigContent } from "@/features/monitors/components/DetailsPage/TcpConfigContent";
import { DnsConfigContent } from "@/features/monitors/components/DetailsPage/DnsConfigContent";
import { Badge } from "@/components/ui/badge.tsx";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "@/lib/consts";
import { UsersApi } from "@/features/users/users-api";
import { Initial } from "@/features/users/components/Initial";
import { StyledLink } from "@/components/common/StyledLink";
import { Flex } from "@/components/common/Flex.tsx";
import { PushConfigContent } from "@/features/monitors/components/DetailsPage/PushConfigContent.tsx";
import { LucideLock } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";

export interface MonitorConfigCardProps {
  monitor: Monitor;
}

const ConfigContent = ({ monitor }: MonitorConfigCardProps) => {
  switch (monitor.type) {
    case "http":
      return <HttpConfigContent monitor={monitor} />;
    case "tcp":
      return <TcpConfigContent monitor={monitor} />;
    case "dns":
      return <DnsConfigContent monitor={monitor} />;
    case "push":
      return <PushConfigContent monitor={monitor} />;
  }
};

const monitorTypeLabelMap: Record<MonitorType, string> = {
  http: "HTTP",
  tcp: "TCP",
  dns: "DNS",
  push: "Push",
};

export const MonitorConfigCard = ({ monitor }: MonitorConfigCardProps) => {
  const { data: users } = useQuery({
    queryKey: [QUERY_KEYS.USERS],
    queryFn: () => UsersApi.getAll(),
  });
  const owner = users?.find((user) => user.id === monitor.ownerId);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Configuration</CardTitle>
        <CardAction>
          <Flex className="items-center">
            <Badge variant="secondary" className="text-sm">
              {monitorTypeLabelMap[monitor.type]}
            </Badge>
            {monitor.source === "config" && (
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
        </CardAction>
      </CardHeader>
      <CardContent>
        <ConfigContent monitor={monitor} />
      </CardContent>
      <CardFooter className="justify-between text-muted-foreground">
        <Flex
          direction="row"
          directionMobile="column"
          className="justify-between w-full"
        >
          <Flex direction="column" directionMobile="column">
            <span>Created {formatDate(new Date(monitor.createdAt))}</span>
            <span>Updated {formatDate(new Date(monitor.updatedAt))}</span>
          </Flex>
          {owner && (
            <span className="flex items-center gap-2">
              Owned by
              <Initial text={owner.username} size="xs" />
              <StyledLink
                to="/user/$username"
                params={{ username: owner.username }}
              >
                {owner.username}
              </StyledLink>
            </span>
          )}
        </Flex>
      </CardFooter>
    </Card>
  );
};
