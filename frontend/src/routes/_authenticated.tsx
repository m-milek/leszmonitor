import { Outlet, createFileRoute } from "@tanstack/react-router";
import { AppSidebar } from "@/components/leszmonitor/sidebar/AppSidebar.tsx";
import { ScrollArea } from "@/components/ui/scroll-area.tsx";

export const Route = createFileRoute("/_authenticated")({
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <div className="flex h-screen w-full text-foreground">
      <AppSidebar />
      <div className="flex-1 flex flex-col bg-background overflow-hidden">
        <ScrollArea className="h-full">
          <main className="min-h-full">
            <Outlet />
          </main>
        </ScrollArea>
      </div>
    </div>
  );
}
