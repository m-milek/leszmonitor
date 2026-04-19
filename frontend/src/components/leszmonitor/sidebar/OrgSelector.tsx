import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useAtomValue } from "jotai";
import { orgAtom } from "@/lib/atoms.ts";
import type { Org } from "@/lib/types.ts";

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

export interface OrgSelectorProps {
  orgs: Org[];
}

export const OrgSelector = ({ orgs }: OrgSelectorProps) => {
  const org = useAtomValue(orgAtom);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild className="w-full h-8">
        <Button variant="secondary">{org?.name}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-full">
        <DropdownMenuGroup>
          {orgs?.map((org) => (
            <OrgEntry key={org.id} orgName={org.name} isCurrent={false} />
          ))}
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
