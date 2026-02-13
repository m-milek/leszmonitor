import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/sidebar/Typography.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useTheme } from "@/components/ui/theme-provider.tsx";

export const Route = createFileRoute("/_authenticated/$teamId/settings/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "light" ? "dark" : "light");
  };

  return (
    <MainPanelContainer>
      <TypographyH1>Settings</TypographyH1>
      <Button onClick={toggleTheme}>Switch Theme</Button>
    </MainPanelContainer>
  );
}
