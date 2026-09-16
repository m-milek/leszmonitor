import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1 } from "@/components/common/Typography.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useTheme } from "@/components/ui/theme-provider.tsx";

export function UserSettingsPage() {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "light" ? "dark" : "light");
  };

  return (
    <PageContainer>
      <TypographyH1>Settings</TypographyH1>
      <Button onClick={toggleTheme}>Switch Theme</Button>
    </PageContainer>
  );
}
