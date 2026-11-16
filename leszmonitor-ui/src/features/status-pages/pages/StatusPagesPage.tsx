import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Card } from "@/components/ui/card";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import { StatusPagesTable } from "@/features/status-pages/components/StatusPagesTable";
import { QUERY_KEYS } from "@/lib/consts";

export function StatusPagesPage() {
  const { data: statusPages } = useQuery({
    queryKey: [QUERY_KEYS.STATUS_PAGES],
    queryFn: () => StatusPagesApi.getAll(),
  });

  return (
    <PageContainer>
      <TypographyH1>Status Pages</TypographyH1>
      <Card>
        {statusPages && <StatusPagesTable statusPages={statusPages} />}
      </Card>
    </PageContainer>
  );
}
