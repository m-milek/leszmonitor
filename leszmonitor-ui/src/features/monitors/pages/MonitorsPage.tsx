import { MonitorsApi } from "@/features/monitors/monitors-api";
import { Link, useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { LucidePlusCircle, LucideX } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { MonitorListItem } from "@/features/monitors/components/MonitorListItem";
import { QUERY_KEYS } from "@/lib/consts";
import { TagsApi } from "@/features/tags/tags-api.ts";
import { useState } from "react";
import { LMInputField } from "@/components/form/LMInputField.tsx";
import { LMFacetedFilter } from "@/components/form/LMFacetedFilter";
import { Tag } from "@/features/tags/components/Tag";
import {
  monitorStatuses,
  type MonitorWithStatus,
} from "@/features/monitors/types.ts";

export function MonitorsPage() {
  const queryClient = useQueryClient();

  const { data: monitors = [] } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS],
    queryFn: () => MonitorsApi.getAll(),
  });

  const { data: tags = [] } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
  });

  const { data: monitorLatestStatuses = [] } = useQuery({
    enabled: !!monitors.length,
    queryKey: [QUERY_KEYS.MONITOR_RESULTS],
    queryFn: () =>
      Promise.all(
        monitors.map((monitor) => MonitorsApi.results.getLatest(monitor.id)),
      ),
  });

  if (!monitors || !tags || !monitorLatestStatuses) {
    return null;
  }

  const { mutateAsync: deleteMutation } = useMutation({
    mutationFn: (monitorId: string) => MonitorsApi.remove(monitorId),
  });

  const navigate = useNavigate();

  const onDeleteMonitor = async (monitorId: string) => {
    await deleteMutation(monitorId);
    queryClient.invalidateQueries({
      queryKey: [QUERY_KEYS.MONITORS],
    });
  };

  const navigateToEditMonitor = (monitorSlug: string) => {
    navigate({
      to: "/monitors/$monitorSlug/edit",
      params: { monitorSlug },
    });
  };

  const monitorsWithStatuses: MonitorWithStatus[] = monitors.map((monitor) => {
    const latestResult = monitorLatestStatuses.find(
      (result) => result?.monitorId === monitor.id,
    );
    return {
      ...monitor,
      status: latestResult?.status ?? "down", // TODO add unknown
    };
  });

  const [filterTags, setFilterTags] = useState<string[]>([]);
  const onTagFilterChange = (selectedTags: string[]) => {
    setFilterTags(selectedTags);
  };

  const [filterName, setFilterName] = useState<string>("");
  const onNameFilterChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFilterName(event.target.value);
  };

  const [statusFilter, setStatusFilter] = useState<string[]>([]);
  const onStatusFilterChange = (selectedStatuses: string[]) => {
    setStatusFilter(selectedStatuses);
  };

  const hasActiveFilters =
    filterName !== "" || filterTags.length > 0 || statusFilter.length > 0;
  const clearFilters = () => {
    setFilterName("");
    setFilterTags([]);
    setStatusFilter([]);
  };

  const filteredMonitors = monitorsWithStatuses.filter((monitor) => {
    const matchesName = monitor.name
      .toLowerCase()
      .includes(filterName.toLowerCase());
    const matchesTags =
      filterTags.length === 0 ||
      filterTags.every((tagId) => monitor.tagIds?.some((id) => id === tagId));
    const matchesStatus =
      statusFilter.length === 0 || statusFilter.includes(monitor.status);

    return matchesName && matchesTags && matchesStatus;
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
                title="Status"
                options={monitorStatuses.map((status) => ({
                  value: status,
                  label: status,
                }))}
                value={statusFilter}
                onChange={onStatusFilterChange}
              />
              {hasActiveFilters && (
                <Button variant="ghost" onClick={clearFilters}>
                  Clear
                  <LucideX />
                </Button>
              )}
            </div>
          </Flex>
        </CardHeader>
        <CardContent>
          <Flex direction="column" className="gap-4">
            {filteredMonitors.map((monitor) => (
              <MonitorListItem
                key={monitor.id}
                monitor={monitor}
                onDeleteMonitor={onDeleteMonitor}
                navigateToEditMonitor={navigateToEditMonitor}
              />
            ))}
          </Flex>
        </CardContent>
      </Card>
    </PageContainer>
  );
}
