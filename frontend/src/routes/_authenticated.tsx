import { Outlet, createFileRoute } from "@tanstack/react-router";
import { AppSidebar } from "@/components/leszmonitor/sidebar/AppSidebar.tsx";

export const Route = createFileRoute("/_authenticated")({
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <div className="flex h-screen w-full">
      <AppSidebar />
      <main className="flex-1 bg-background">
        <Outlet />
      </main>
    </div>
  );
}
