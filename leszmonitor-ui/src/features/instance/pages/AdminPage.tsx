import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { UsersCard } from "@/features/instance/components/UsersCard";
import { GlobalParametersConfiguration } from "@/features/instance/components/GlobalParametersConfiguration.tsx";

export function AdminPage() {
  return (
    <PageContainer>
      <TypographyH1>Administration Dashboard</TypographyH1>
      <UsersCard />
      <GlobalParametersConfiguration />
    </PageContainer>
  );
}
