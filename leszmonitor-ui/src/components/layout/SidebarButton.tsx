import {
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import { cn } from "@/lib/utils.ts";
import { Link, useLocation } from "@tanstack/react-router";

interface SidebarButtonProps {
  icon: React.ReactNode;
  href: string;
  label: string;
}

export const SidebarButton = ({ icon, href, label }: SidebarButtonProps) => {
  const location = useLocation();
  const matchesCurrentUrl = location.pathname === href;

  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        asChild
        isActive={matchesCurrentUrl}
        className={cn("data-[active=true]:text-sidebar-primary")}
      >
        <Link to={href} draggable={false}>
          {icon}
          <span>{label}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
};
