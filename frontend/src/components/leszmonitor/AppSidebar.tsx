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
import { Logo } from "@/components/leszmonitor/Logo.tsx";
import { useAtom } from "jotai";
import { userAtom, usernameAtom } from "@/lib/atoms.ts";
import { useQuery } from "@tanstack/react-query";
import { fetchUser } from "@/lib/fetchUser.ts";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";

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
  const [username, setUsername] = useAtom(usernameAtom);
  const [user, setUser] = useAtom(userAtom);
  const hasInitialized = useRef(false);

  useEffect(() => {
    if (hasInitialized.current) return;
    hasInitialized.current = true;

    const getTokenAndExtractUsername = async () => {
      const token = await cookieStore.get("LOGIN_TOKEN");
      if (token?.value) {
        const claims = jwtDecode(token.value) as JwtClaims;
        if (claims?.username) {
          setUsername(claims.username);
        }
      }
    };

    getTokenAndExtractUsername();
  }, [setUsername]);

  const { data: userData } = useQuery({
    queryKey: ["user", username],
    queryFn: () => fetchUser(username!),
    enabled: !!username,
    staleTime: 5 * 60 * 1000,
  });

  useEffect(() => {
    if (userData) {
      setUser(userData);
    }
  }, [userData, setUser]);

  if (!username) {
    return null;
  }

  return (
    <Sidebar variant="inset">
      <SidebarHeader>
        <div className="p-2">
          <Logo />
        </div>
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
                label="New Monitor"
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
              <SidebarButton
                icon={<LucideSearch />}
                href="/search"
                label="Search"
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
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        {user ? (
          <div className="flex items-center space-x-2">
            <div className="h-8 w-8 rounded-full bg-muted" />
            <span>{user.username}</span>
          </div>
        ) : (
          <span>Not logged in</span>
        )}
      </SidebarFooter>
    </Sidebar>
  );
};
