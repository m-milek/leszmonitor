import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AppSidebar } from "@/components/layout/AppSidebar";
import { WebSocketProvider } from "@/app/providers/WebSocketProvider";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { getCookie } from "@/lib/cookies";

export const Route = createFileRoute("/_authenticated")({
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <WebSocketProvider>
      <SidebarProvider
        defaultOpen={getCookie("sidebar_state") !== "false"}
        style={{ "--sidebar-width": "16rem" } as React.CSSProperties}
      >
        <AppSidebar />
        <SidebarInset>
          <header className="flex h-12 shrink-0 items-center px-4 md:hidden">
            <SidebarTrigger />
          </header>
          <Outlet />
        </SidebarInset>
      </SidebarProvider>
    </WebSocketProvider>
  );
}
