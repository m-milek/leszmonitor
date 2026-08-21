import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import {
  TypographyH1,
  TypographyH3,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { Center } from "./ui/Center.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

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
