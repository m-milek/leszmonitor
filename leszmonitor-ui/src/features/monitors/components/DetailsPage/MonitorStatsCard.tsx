import type { MonitorStats } from "@/features/monitors/types.ts";
import { Card } from "@/components/ui/card.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";

const MS_SUFFIX = " ms";

export interface MonitorStatsCardProps {
  stats?: MonitorStats;
}

export function MonitorStatsCard({ stats }: MonitorStatsCardProps) {
  return (
    <Card>
      <Flex direction="row" className="gap-2">
        <MonitorStatsCardItemWrapper title="Uptime">
          <MonitorStatsNumber
            value={stats && (stats.uptime?.statusToPercentage?.["up"] ?? 0)}
            decimalPlaces={2}
            suffix={"%"}
          />
        </MonitorStatsCardItemWrapper>
        <Separator orientation="vertical" />
        <MonitorStatsCardItemWrapper title="Average Latency">
          <MonitorStatsNumber value={stats?.latency?.avg} suffix={MS_SUFFIX} />
        </MonitorStatsCardItemWrapper>
        <Separator orientation="vertical" />
        <MonitorStatsCardItemWrapper title="Minimum Latency">
          <MonitorStatsNumber value={stats?.latency?.min} suffix={MS_SUFFIX} />
        </MonitorStatsCardItemWrapper>
        <Separator orientation="vertical" />
        <MonitorStatsCardItemWrapper title="Maximum Latency">
          <MonitorStatsNumber value={stats?.latency?.max} suffix={MS_SUFFIX} />
        </MonitorStatsCardItemWrapper>
      </Flex>
    </Card>
  );
}

const MonitorStatsCardItemWrapper = ({
  children,
  title,
}: {
  children: React.ReactNode;
  title: string;
}) => {
  return (
    <Flex
      direction="column"
      className="flex-1 items-center justify-center gap-1"
    >
      {children}
      <div className="text-muted-foreground">{title}</div>
    </Flex>
  );
};

interface MonitorStatsNumberProps {
  value?: number;
  decimalPlaces?: number;
  suffix?: string;
}

const MonitorStatsNumber = ({
  value,
  decimalPlaces = 0,
  suffix,
}: MonitorStatsNumberProps) => {
  if (value === undefined) {
    return <Skeleton className="h-8 w-24" />;
  }

  const formattedValue = value.toLocaleString(undefined, {
    minimumFractionDigits: decimalPlaces,
    maximumFractionDigits: decimalPlaces,
  });

  return (
    <span className="text-2xl font-semibold">
      {formattedValue}
      {suffix && <span>{suffix}</span>}
    </span>
  );
};
