import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1, TypographyH3 } from "@/components/common/Typography";
import { Center } from "@/components/common/Center";
import { Flex } from "@/components/common/Flex";

export function GlobalNotFound() {
  return (
    <PageContainer className="flex h-screen w-full bg-background color-text">
      <Center>
        <Flex direction="column" className="gap-4 items-center">
          <TypographyH1>404</TypographyH1>
          <TypographyH3>Not Found</TypographyH3>
        </Flex>
      </Center>
    </PageContainer>
  );
}
