import { SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { Link, useLocation } from "@tanstack/react-router";

interface SidebarButtonProps {
  icon: React.ReactNode;
  href: string;
  label: string;
  /** Served by the Go backend, not the router — needs a full page load. */
  external?: boolean;
}

export const SidebarButton = ({
  icon,
  href,
  label,
  external,
}: SidebarButtonProps) => {
  const location = useLocation();
  const matchesCurrentUrl = location.pathname === href;

  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        render={
          external ? (
            <a href={href} draggable={false} />
          ) : (
            <Link to={href} draggable={false} />
          )
        }
        isActive={matchesCurrentUrl}
        className="transition-colors hover:text-sidebar-foreground active:translate-y-px data-active:text-sidebar-primary [&[data-active]:hover]:text-sidebar-primary"
      >
        {icon}
        <span>{label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
};
