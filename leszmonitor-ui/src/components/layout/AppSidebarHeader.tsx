import { SidebarHeader } from "@/components/ui/sidebar";
import { Flex } from "@/components/common/Flex";
import { Link } from "@tanstack/react-router";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo";
import { WebSocketStatusIndicator } from "@/components/layout/WebSocketStatusIndicator";

export function AppSidebarHeader() {
  return (
    <SidebarHeader className="p-2">
      <Flex direction="row" className="justify-between items-center">
        <div className="p-2">
          <Link to={"/monitors"}>
            <LeszmonitorLogo />
          </Link>
        </div>
        <div className="p-2">
          <WebSocketStatusIndicator />
        </div>
      </Flex>
    </SidebarHeader>
  );
}
