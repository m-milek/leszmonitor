import { SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
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
        render={<Link to={href} draggable={false} />}
        isActive={matchesCurrentUrl}
        className="data-active:text-sidebar-primary"
      >
        {icon}
        <span>{label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
};
