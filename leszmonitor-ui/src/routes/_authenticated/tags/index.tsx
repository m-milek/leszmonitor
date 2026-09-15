import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { TagsTable } from "@/components/leszmonitor/tables/TagsTable.tsx";
import { NewTagDialog } from "@/components/leszmonitor/dialogs/NewTagDialog.tsx";
import { getAllTags } from "@/lib/data/tags-api.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";

export const Route = createFileRoute("/_authenticated/tags/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: tags } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => getAllTags(),
  });

  return (
    <PageContainer>
      <Flex className="items-center justify-between">
        <TypographyH1>Tags</TypographyH1>
        <NewTagDialog />
      </Flex>
      {tags ? <TagsTable tags={tags} /> : <Skeleton className="h-32" />}
    </PageContainer>
  );
}
