import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AppSidebar } from "@/components/leszmonitor/sidebar/AppSidebar.tsx";
import { WebSocketProvider } from "@/components/leszmonitor/providers/WebSocketProvider.tsx";
import { ScrollArea } from "@/components/ui/scroll-area.tsx";

export const Route = createFileRoute("/_authenticated")({
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <WebSocketProvider>
      <div className="grid min-h-svh w-full text-foreground md:grid-cols-[18rem_minmax(0,1fr)]">
        <AppSidebar />
        <main className="min-w-0 bg-background">
          <ScrollArea className="h-svh">
            <Outlet />
          </ScrollArea>
        </main>
      </div>
    </WebSocketProvider>
  );
}
