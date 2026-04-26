import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useAppStore } from "@/lib/store.ts";
import type { Project } from "@/lib/types.ts";

interface ProjectEntryProps {
  projectName: string;
  isCurrent: boolean;
}
export const ProjectEntry = ({ projectName, isCurrent }: ProjectEntryProps) => {
  return (
    <DropdownMenuItem>
      {projectName} {isCurrent && "(current)"}
    </DropdownMenuItem>
  );
};

export interface ProjectSelectorProps {
  projects: Project[];
}

export const ProjectSelector = ({ projects }: ProjectSelectorProps) => {
  const { project } = useAppStore();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild className="w-full h-8">
        <Button variant="secondary">{project?.name}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-full">
        <DropdownMenuGroup>
          {projects?.map((p) => (
            <ProjectEntry key={p.id} projectName={p.name} isCurrent={false} />
          ))}
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
