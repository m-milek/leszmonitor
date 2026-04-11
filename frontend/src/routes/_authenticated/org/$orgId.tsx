import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useSetAtom } from "jotai";
import { orgAtom } from "@/lib/atoms.ts";
import { useEffect } from "react";
import { getOrg } from "@/lib/data/orgData.ts";
import { useQuery } from "@tanstack/react-query";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";

export const Route = createFileRoute("/_authenticated/org/$orgId")({
  component: OrgLayout,
  notFoundComponent: NotFound,
});

function OrgLayout() {
  const { orgId } = Route.useParams();
  const setOrgAtom = useSetAtom(orgAtom);

  const { data: org } = useQuery({
    queryKey: ["org", orgId],
    queryFn: () => getOrg(orgId),
  });

  useEffect(() => {
    if (org) {
      setOrgAtom(org);
    }
  }, [org, orgId, setOrgAtom]);

  return <Outlet />;
}

function NotFound() {
  return (
    <MainPanelContainer>
      <TypographyH1>Not Found</TypographyH1>;
    </MainPanelContainer>
  );
}
