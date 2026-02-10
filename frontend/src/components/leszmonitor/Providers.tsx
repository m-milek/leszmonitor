import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Provider } from "jotai";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { store } from "@/lib/atoms.ts";
import { ThemeProvider } from "@/components/ui/theme-provider.tsx";

const queryClient = new QueryClient();

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <Provider store={store}>
          <SidebarProvider>
            {children}
            <TanStackRouterDevtools position="bottom-right" />
          </SidebarProvider>
        </Provider>
      </QueryClientProvider>
    </ThemeProvider>
  );
};
