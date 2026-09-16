import { Link, useNavigate } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { LucidePlusCircle } from "lucide-react";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1 } from "@/components/common/Typography.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  deleteMonitor,
  getAllMonitors,
} from "@/features/monitors/api/monitors.ts";
import { MonitorListItem } from "@/features/monitors/components/MonitorListItem.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";

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
