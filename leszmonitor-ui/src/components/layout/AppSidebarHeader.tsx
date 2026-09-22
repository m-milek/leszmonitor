import { SidebarHeader } from "@/components/ui/sidebar";
import { Flex } from "@/components/common/Flex";
import { Link } from "@tanstack/react-router";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo";
import { WebSocketStatusIndicator } from "@/components/layout/WebSocketStatusIndicator";

export function AppSidebarHeader() {
  return (
    <SidebarHeader>
      <Flex direction="row" className="justify-between items-center">
        <Link to={"/monitors"}>
          <LeszmonitorLogo />
        </Link>
        <WebSocketStatusIndicator />
      </Flex>
    </SidebarHeader>
  );
}
