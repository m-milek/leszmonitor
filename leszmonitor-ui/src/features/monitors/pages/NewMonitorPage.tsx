import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { NewMonitorForm } from "@/features/monitors/forms/NewMonitorForm";

export function NewMonitorPage() {
  return (
    <PageContainer>
      <TypographyH1>New Monitor Wizard</TypographyH1>
      <Card>
        <CardContent>
          <NewMonitorForm formId="new-monitor-form" />
        </CardContent>
        <CardFooter>
          <Button type="submit" form="new-monitor-form">
            Create Monitor
          </Button>
        </CardFooter>
      </Card>
    </PageContainer>
  );
}
