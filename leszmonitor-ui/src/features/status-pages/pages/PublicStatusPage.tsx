import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Center } from "@/components/common/Center";
import { Flex } from "@/components/common/Flex";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo";
import { StatusDot } from "@/components/common/StatusDot";
import { TypographyH1, TypographyH3 } from "@/components/common/Typography";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Tag } from "@/features/tags/components/Tag";
import { PublicMonitorRow } from "@/features/status-pages/components/PublicMonitorRow";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import type { PublicMonitor } from "@/features/status-pages/types";
import type { Tag as TagModel } from "@/features/tags/types";
import { QUERY_KEYS } from "@/lib/consts";

const REFETCH_INTERVAL_MS = 30_000;

export interface PublicStatusPageProps {
  slug: string;
}

interface MonitorGroup {
  tag?: TagModel;
  monitors: PublicMonitor[];
}

const groupByTag = (monitors: PublicMonitor[]): MonitorGroup[] => {
  const groups = new Map<string, MonitorGroup>();
  for (const monitor of monitors) {
    for (const tag of monitor.tags) {
      const group = groups.get(tag.id) ?? { tag, monitors: [] };
      group.monitors.push(monitor);
      groups.set(tag.id, group);
    }
  }
  return [...groups.values()].sort((a, b) =>
    (a.tag?.name ?? "").localeCompare(b.tag?.name ?? ""),
  );
};

const OverallStatus = ({ monitors }: { monitors: PublicMonitor[] }) => {
  const downCount = monitors.filter((m) => m.status === "down").length;
  const isAllUp = monitors.every((m) => m.status === "up");

  if (downCount > 0) {
    return (
      <Flex direction="row" className="items-center gap-3">
        <StatusDot status="down" size="lg" />
        <TypographyH3>
          {downCount} of {monitors.length} services down
        </TypographyH3>
      </Flex>
    );
  }

  return (
    <Flex direction="row" className="items-center gap-3">
      <StatusDot status={isAllUp ? "up" : "unknown"} size="lg" />
      <TypographyH3>
        {isAllUp ? "All systems operational" : "Some services have no data"}
      </TypographyH3>
    </Flex>
  );
};

const PageMessage = ({
  title,
  children,
}: {
  title: string;
  children?: React.ReactNode;
}) => (
  <Center className="h-svh">
    <Flex direction="column" className="items-center gap-4">
      <TypographyH1>{title}</TypographyH1>
      {children}
    </Flex>
  </Center>
);

export function PublicStatusPage({ slug }: PublicStatusPageProps) {
  const pageQuery = useQuery({
    queryKey: [QUERY_KEYS.PUBLIC_STATUS_PAGES, slug],
    queryFn: () => StatusPagesApi.public.get(slug),
    retry: false,
  });

  const monitorsQuery = useQuery({
    queryKey: [QUERY_KEYS.PUBLIC_STATUS_PAGE_MONITORS, slug],
    queryFn: () => StatusPagesApi.public.getMonitors(slug),
    enabled: pageQuery.isSuccess,
    refetchInterval: REFETCH_INTERVAL_MS,
  });

  if (pageQuery.error) {
    const status = pageQuery.error.cause;

    if (status === 401 || status === 403) {
      return (
        <PageMessage title="Login required">
          <Link to="/login" className="underline">
            Log in to view this status page
          </Link>
        </PageMessage>
      );
    }
    if (status === 404) {
      return <PageMessage title="Status page not found" />;
    }
    return <PageMessage title="Failed to load status page" />;
  }

  const page = pageQuery.data;
  const monitors = monitorsQuery.data;
  const groups: MonitorGroup[] = monitors
    ? page?.groupMonitorsByTag
      ? groupByTag(monitors)
      : [{ monitors }]
    : [];

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col gap-6 p-6 max-md:p-4">
      <LeszmonitorLogo />
      {page ? (
        <Flex direction="column" className="gap-1">
          <TypographyH1>{page.title}</TypographyH1>
          {page.description && (
            <span className="text-muted-foreground">{page.description}</span>
          )}
        </Flex>
      ) : (
        <Skeleton className="h-9 w-64" />
      )}
      {monitors ? (
        <>
          <Card>
            <CardContent>
              <OverallStatus monitors={monitors} />
            </CardContent>
          </Card>
          {groups.map((group) => (
            <Card key={group.tag?.id ?? "all"}>
              <CardContent>
                {group.tag && <Tag tag={group.tag} />}
                <Flex direction="column" className="divide-y">
                  {group.monitors.map((monitor) => (
                    <PublicMonitorRow key={monitor.id} monitor={monitor} />
                  ))}
                </Flex>
              </CardContent>
            </Card>
          ))}
        </>
      ) : (
        <Skeleton className="h-64 w-full" />
      )}
    </div>
  );
}
