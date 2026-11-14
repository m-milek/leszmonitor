import { MonitorsApi } from "@/features/monitors/monitors-api";
import { Link } from "@tanstack/react-router";
import { useQueries, useQuery } from "@tanstack/react-query";
import { LucidePlusCircle, LucideX } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { getEffectiveStatus } from "@/features/monitors/status";
import { MonitorStatusSummary } from "@/features/monitors/components/MonitorStatusSummary";
import { MonitorsTable } from "@/features/monitors/components/MonitorsTable";
import { QUERY_KEYS } from "@/lib/consts";
import { TagsApi } from "@/features/tags/tags-api.ts";
import { LMInputField } from "@/components/form/LMInputField.tsx";
import { LMFacetedFilter } from "@/components/form/LMFacetedFilter";
import { Tag } from "@/features/tags/components/Tag";
import {
  monitorTypes,
  type MonitorStatus,
  type MonitorType,
  type MonitorWithStatus,
} from "@/features/monitors/types.ts";
import type { MonitorsSearch } from "@/features/monitors/monitors-search";

export interface MonitorsPageProps {
  search: MonitorsSearch;
  onSearchChange: (search: MonitorsSearch) => void;
}

export function MonitorsPage({
  search,
  onSearchChange,
}: Readonly<MonitorsPageProps>) {
  const { data: monitors = [], isPending: isMonitorsPending } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS],
    queryFn: () => MonitorsApi.getAll(),
  });

  const { data: tags = [] } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
  });

  const latestResultQueries = useQueries({
    queries: monitors.map((monitor) => ({
      queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id, "latest"],
      queryFn: () => MonitorsApi.results.getLatest(monitor.id),
    })),
  });

  const monitorsWithStatuses: MonitorWithStatus[] = monitors.map(
    (monitor, index) => ({
      ...monitor,
      status: latestResultQueries[index]?.data?.status ?? "unknown",
    }),
  );

  const filterTags = search.tags ?? [];
  const onTagFilterChange = (selectedTags: string[]) => {
    onSearchChange({
      ...search,
      tags: selectedTags.length > 0 ? selectedTags : undefined,
    });
  };

  const filterName = search.q ?? "";
  const onNameFilterChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onSearchChange({ ...search, q: event.target.value || undefined });
  };

  const statusFilter: string[] = search.status ?? [];
  const onStatusFilterChange = (selectedStatuses: string[]) => {
    onSearchChange({
      ...search,
      status:
        selectedStatuses.length > 0
          ? (selectedStatuses as MonitorStatus[])
          : undefined,
    });
  };

  const typeFilter: string[] = search.type ?? [];
  const onTypeFilterChange = (selectedTypes: string[]) => {
    onSearchChange({
      ...search,
      type:
        selectedTypes.length > 0 ? (selectedTypes as MonitorType[]) : undefined,
    });
  };

  const hasActiveFilters =
    filterName !== "" ||
    filterTags.length > 0 ||
    statusFilter.length > 0 ||
    typeFilter.length > 0;
  const clearFilters = () => {
    onSearchChange({});
  };

  const filteredMonitors = monitorsWithStatuses.filter((monitor) => {
    const matchesName = monitor.name
      .toLowerCase()
      .includes(filterName.toLowerCase());
    const matchesTags =
      filterTags.length === 0 ||
      filterTags.every((tagId) => monitor.tagIds?.some((id) => id === tagId));
    const matchesStatus =
      statusFilter.length === 0 ||
      statusFilter.includes(getEffectiveStatus(monitor));

    const matchesType =
      typeFilter.length === 0 || typeFilter.includes(monitor.type);

    return matchesName && matchesTags && matchesStatus && matchesType;
  });

  return (
    <PageContainer>
      <TypographyH1>Monitors</TypographyH1>
      <Card>
        <CardHeader>
          <Flex direction="column" className="gap-2">
            <Link to={"/monitors/new"} className="mr-auto">
              <Button>
                <LucidePlusCircle />
                <span>New Monitor</span>
              </Button>
            </Link>
            <div className="flex flex-wrap items-center gap-2">
              <div className="w-64">
                <LMInputField
                  placeholder={"Search..."}
                  name="monitor-name-filter"
                  value={filterName}
                  onChange={onNameFilterChange}
                />
              </div>
              <LMFacetedFilter
                title="Tags"
                options={tags.map((tag) => ({
                  value: tag.id,
                  label: tag.name,
                  render: <Tag tag={tag} />,
                }))}
                value={filterTags}
                onChange={onTagFilterChange}
              />
              <LMFacetedFilter
                title="Type"
                options={monitorTypes.map((type) => ({
                  value: type,
                  label: type.toUpperCase(),
                }))}
                value={typeFilter}
                onChange={onTypeFilterChange}
              />
              {hasActiveFilters && (
                <Button variant="ghost" onClick={clearFilters}>
                  Clear
                  <LucideX />
                </Button>
              )}
              <MonitorStatusSummary
                className="ml-auto"
                monitors={monitorsWithStatuses}
                value={statusFilter}
                onChange={onStatusFilterChange}
              />
            </div>
          </Flex>
        </CardHeader>
        <CardContent>
          {isMonitorsPending ? (
            <Flex direction="column" className="gap-2">
              {Array.from({ length: 5 }, (_, i) => (
                <Skeleton key={i} className="h-9 w-full" />
              ))}
            </Flex>
          ) : (
            <MonitorsTable monitors={filteredMonitors} tags={tags} />
          )}
        </CardContent>
      </Card>
    </PageContainer>
  );
}
