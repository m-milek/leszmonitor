import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Provider as JotaiProvider } from "jotai";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { store } from "@/lib/atoms.ts";
import { ThemeProvider } from "@/components/ui/theme-provider.tsx";
import { TooltipProvider } from "@/components/ui/tooltip.tsx";

const queryClient = new QueryClient();

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <JotaiProvider store={store}>
          <TooltipProvider>
            <SidebarProvider>
              {children}
              <TanStackRouterDevtools position="bottom-right" />
            </SidebarProvider>
          </TooltipProvider>
        </JotaiProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
};
