import { PageContainer } from "@/components/common/PageContainer.tsx";
import {
  TypographyH1,
  TypographyH3,
} from "@/components/common/Typography.tsx";
import { Center } from "@/components/common/Center.tsx";
import { Flex } from "@/components/common/Flex.tsx";

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
