import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import {
  LayoutDashboardIcon,
  LucideBookText,
  LucideGroup,
  LucidePlusCircle,
  LucideScanEye,
  LucideSearch,
  LucideSettings,
  LucideUsers,
} from "lucide-react";

interface SidebarButtonProps {
  icon: React.ReactNode;
  href: string;
  label: string;
  variant?: "default" | "primary";
}

const SidebarButton = ({
  icon,
  href,
  label,
  variant = "default",
}: SidebarButtonProps) => {
  const onClick = () => {
    console.log("navigating to", href);
  };

  const className =
    variant === "primary"
      ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
      : "";

  return (
    <SidebarMenuItem>
      <SidebarMenuButton onClick={onClick} className={className}>
        {icon}
        <span>{label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
};

export const AppSidebar = () => {
  return (
    <Sidebar>
      <SidebarHeader>
        Leszmonitor
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton>Select Team</SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem>
                  <span>Acme Inc</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucidePlusCircle />}
                href="/new-monitor"
                label="Quick Create"
                variant="primary"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LayoutDashboardIcon />}
                href="/dashboard"
                label="Dashboard"
              />
              <SidebarButton
                icon={<LucideScanEye />}
                href="/monitors"
                label="Monitors"
              />
              <SidebarButton
                icon={<LucideGroup />}
                href="/groups"
                label="Groups"
              />
              <SidebarButton icon={<LucideUsers />} href="/team" label="Team" />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucideBookText />}
                href="/documentation"
                label="Documentation"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarButton
              icon={<LucideSettings />}
              href="/settings"
              label="Settings"
            />
            <SidebarButton
              icon={<LucideSearch />}
              href="/search"
              label="Search"
            />
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>Footer</SidebarFooter>
    </Sidebar>
  );
};
