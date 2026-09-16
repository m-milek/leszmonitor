import { lazy, Suspense } from "react";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { ThemeProvider } from "@/components/ui/theme-provider.tsx";
import { TooltipProvider } from "@/components/ui/tooltip.tsx";

// Keep the router devtools out of production bundles. `import.meta.env.PROD` is
// statically inlined by Vite, so the dynamic import below is dropped entirely
// from prod builds.
const TanStackRouterDevtools = import.meta.env.PROD
  ? () => null
  : lazy(() =>
      import("@tanstack/react-router-devtools").then((m) => ({
        default: m.TanStackRouterDevtools,
      })),
    );

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider>
      <TooltipProvider>
        <SidebarProvider>
          {children}
          <Suspense fallback={null}>
            <TanStackRouterDevtools position="bottom-right" />
          </Suspense>
        </SidebarProvider>
      </TooltipProvider>
    </ThemeProvider>
  );
};
