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
import { LucideDot } from "lucide-react";
import { Badge } from "@/components/ui/badge.tsx";

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
  }
};

const monitorTypeLabelMap: Record<MonitorType, string> = {
  http: "HTTP",
  tcp: "TCP",
  dns: "DNS",
};

export const MonitorConfigCard = ({ monitor }: MonitorConfigCardProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Configuration</CardTitle>
        <CardAction>
          <Badge variant="secondary" className="text-sm">
            {monitorTypeLabelMap[monitor.type]}
          </Badge>
        </CardAction>
      </CardHeader>
      <CardContent>
        <ConfigContent monitor={monitor} />
      </CardContent>
      <CardFooter className="text-muted-foreground">
        Created {formatDate(new Date(monitor.createdAt))} <LucideDot /> Updated{" "}
        {formatDate(new Date(monitor.updatedAt))}
      </CardFooter>
    </Card>
  );
};
