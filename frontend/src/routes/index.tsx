import { createFileRoute } from "@tanstack/react-router";
import { AppSidebar } from "@/components/leszmonitor/sidebar/AppSidebar.tsx";

export const Route = createFileRoute("/")({
  component: App,
});

function App() {
  return (
    <div className="flex h-screen w-full">
      <AppSidebar />
      <main className="flex-1 bg-background">main</main>
    </div>
  );
}
