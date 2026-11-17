import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { UsersCard } from "@/features/settings/components/UsersCard";
import { GlobalParametersConfiguration } from "@/features/settings/components/GlobalParametersConfiguration.tsx";
import { NotificationChannelsConfiguration } from "@/features/settings/components/NotificationChannelsConfiguration..tsx";

export function AdminPage() {
  return (
    <PageContainer>
      <TypographyH1>Administration Dashboard</TypographyH1>
      <UsersCard />
      <GlobalParametersConfiguration />
      <NotificationChannelsConfiguration />
    </PageContainer>
  );
}
