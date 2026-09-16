import { TagsApi } from "@/features/tags/tags-api";
import { useQuery } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import { Skeleton } from "@/components/ui/skeleton";
import { TagsTable } from "@/features/tags/components/TagsTable";
import { NewTagDialog } from "@/features/tags/components/NewTagDialog";
import { QUERY_KEYS } from "@/lib/consts";

export function TagsPage() {
  const { data: tags } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
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
