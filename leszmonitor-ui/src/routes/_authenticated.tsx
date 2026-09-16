import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AppSidebar } from "@/components/layout/AppSidebar";
import { WebSocketProvider } from "@/app/providers/WebSocketProvider";
import { ScrollArea } from "@/components/ui/scroll-area";

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
