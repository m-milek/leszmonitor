import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/components/ui/theme-provider.tsx";
import { TooltipProvider } from "@/components/ui/tooltip.tsx";

const queryClient = new QueryClient();

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <TooltipProvider>
          <SidebarProvider>
            {children}
            <TanStackRouterDevtools position="bottom-right" />
          </SidebarProvider>
        </TooltipProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
};
