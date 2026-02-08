import { createFileRoute } from "@tanstack/react-router";
import { AppSidebar } from "@/components/leszmonitor/AppSidebar.tsx";

export const Route = createFileRoute("/")({
  component: App,
});

function App() {
  return (
    <div className="flex">
      <AppSidebar />
      <main>main</main>
    </div>
  );
}
