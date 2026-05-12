import { SidebarHeader } from "@/components/ui/sidebar.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Link } from "@tanstack/react-router";
import { LeszmonitorLogo } from "@/components/leszmonitor/ui/LeszmonitorLogo.tsx";
import { WebSocketStatusIndicator } from "@/components/leszmonitor/sidebar/WebSocketStatusIndicator.tsx";

export function AppSidebarHeader() {
  return (
    <SidebarHeader>
      <Flex direction="row" className="justify-between items-center">
        <div className="p-2">
          <Link to={"/projects"}>
            <LeszmonitorLogo />
          </Link>
        </div>
        <WebSocketStatusIndicator />
      </Flex>
    </SidebarHeader>
  );
}
