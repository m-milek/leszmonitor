import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { ThemeProvider } from "@/components/ui/theme-provider.tsx";
import { TooltipProvider } from "@/components/ui/tooltip.tsx";

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider>
      <TooltipProvider>
        <SidebarProvider>
          {children}
          <TanStackRouterDevtools position="bottom-right" />
        </SidebarProvider>
      </TooltipProvider>
    </ThemeProvider>
  );
};
