import { Link } from "@tanstack/react-router";
import { LucidePlusCircle } from "lucide-react";
import { Flex } from "@/components/common/Flex";
import { Button } from "@/components/ui/button";

export interface MonitorsEmptyStateProps {
  hasMonitors: boolean;
  onClearFilters: () => void;
}

export const MonitorsEmptyState = ({
  hasMonitors,
  onClearFilters,
}: MonitorsEmptyStateProps) => {
  if (!hasMonitors) {
    return (
      <Flex direction="column" className="items-center gap-2">
        <span>No monitors yet.</span>
        <Link to="/monitors/new">
          <Button>
            <LucidePlusCircle />
            <span>Create your first monitor</span>
          </Button>
        </Link>
      </Flex>
    );
  }

  return (
    <Flex direction="column" className="items-center gap-2">
      <span>No monitors match the filters.</span>
      <Button variant="outline" onClick={onClearFilters}>
        Clear filters
      </Button>
    </Flex>
  );
};
