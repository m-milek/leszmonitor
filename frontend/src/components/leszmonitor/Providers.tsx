import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Provider } from "jotai";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <Provider>
        <SidebarProvider>
          {children}
          <TanStackRouterDevtools position="bottom-right" />
        </SidebarProvider>
      </Provider>
    </QueryClientProvider>
  );
};
