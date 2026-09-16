import { Link, useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { LucidePlusCircle } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  deleteMonitor,
  getAllMonitors,
} from "@/features/monitors/monitors-api";
import { MonitorListItem } from "@/features/monitors/components/MonitorListItem";
import { QUERY_KEYS } from "@/lib/consts";

export function MonitorsPage() {
  const queryClient = useQueryClient();

  const { data: monitors = [] } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS],
    queryFn: () => getAllMonitors(),
  });

  const { mutateAsync: deleteMutation } = useMutation({
    mutationFn: (monitorId: string) => deleteMonitor(monitorId),
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

  return (
    <PageContainer>
      <TypographyH1>Monitors</TypographyH1>
      <Card>
        <CardHeader>
          <Link to={"/monitors/new"}>
            <Button>
              <LucidePlusCircle />
              <span>New Monitor</span>
            </Button>
          </Link>
        </CardHeader>
        <CardContent>
          <Flex direction="column" className="gap-4">
            {monitors.map((monitor) => (
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
