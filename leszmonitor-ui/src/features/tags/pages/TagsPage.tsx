import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1 } from "@/components/common/Typography.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { TagsTable } from "@/features/tags/components/TagsTable.tsx";
import { NewTagDialog } from "@/features/tags/components/NewTagDialog.tsx";
import { getAllTags } from "@/features/tags/api/tags.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";

export function TagsPage() {
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
