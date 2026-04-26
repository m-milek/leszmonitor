import {
  CircleCheckIcon,
  InfoIcon,
  Loader2Icon,
  TriangleAlertIcon,
  XCircleIcon,
} from "lucide-react";
import { useTheme } from "next-themes";
import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme();

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      className="toaster group"
      toastOptions={{
        classNames: {
          icon: "size-6 pr-5",
        },
      }}
      icons={{
        success: <CircleCheckIcon className="text-green-500" />,
        info: <InfoIcon className="text-blue-500 " />,
        warning: <TriangleAlertIcon className="text-yellow-500" />,
        error: <XCircleIcon className="text-destructive" />,
        loading: <Loader2Icon className="animate-spin text-foreground" />,
      }}
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
          "--border-radius": "var(--radius)",
        } as React.CSSProperties
      }
      {...props}
    />
  );
};

export { Toaster };
