import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import {
  TypographyH1,
  TypographyH3,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { Center } from "@/components/leszmonitor/ui/Center.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

export function GlobalNotFound() {
  return (
    <div className="flex h-screen w-full">
      <MainPanelContainer>
        <Center>
          <Flex direction="vertical" align="center" gap={1}>
            <TypographyH1>404</TypographyH1>
            <TypographyH3>Not Found</TypographyH3>
          </Flex>
        </Center>
      </MainPanelContainer>
    </div>
  );
}
