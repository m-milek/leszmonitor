import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useQuery } from "@tanstack/react-query";
import { useAtomValue } from "jotai";
import { orgAtom } from "@/lib/atoms.ts";
import { fetchOrgs } from "@/lib/data/orgData.ts";

interface OrgEntryProps {
  orgName: string;
  isCurrent: boolean;
}
export const OrgEntry = ({ orgName, isCurrent }: OrgEntryProps) => {
  return (
    <DropdownMenuItem>
      {orgName} {isCurrent && "(current)"}
    </DropdownMenuItem>
  );
};

export const OrgSelector = () => {
  const { data: orgsData } = useQuery({
    queryKey: ["orgs"],
    queryFn: fetchOrgs,
  });

  const org = useAtomValue(orgAtom);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild className="w-full">
        <Button variant="secondary">{org?.name}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-full">
        <DropdownMenuGroup>
          {orgsData?.map((orgItem) => (
            <OrgEntry
              key={orgItem.id}
              orgName={orgItem.name}
              isCurrent={false}
            />
          ))}
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

