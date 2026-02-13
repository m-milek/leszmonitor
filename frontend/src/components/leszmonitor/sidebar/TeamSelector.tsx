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
import { teamAtom } from "@/lib/atoms.ts";
import { fetchTeams } from "@/lib/data/teamData.ts";

interface TeamEntryProps {
  teamName: string;
  isCurrent: boolean;
}
export const TeamEntry = ({ teamName, isCurrent }: TeamEntryProps) => {
  return (
    <DropdownMenuItem>
      {teamName} {isCurrent && "(current)"}
    </DropdownMenuItem>
  );
};

export const TeamSelector = () => {
  const { data: teamsData } = useQuery({
    queryKey: ["teams"],
    queryFn: fetchTeams,
  });

  const team = useAtomValue(teamAtom);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild className="w-full">
        <Button variant="secondary">{team?.name}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-full">
        <DropdownMenuGroup>
          {teamsData?.map((team) => (
            <TeamEntry key={team.id} teamName={team.name} isCurrent={false} />
          ))}
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
