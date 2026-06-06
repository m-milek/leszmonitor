import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

export interface ResourceDiffProps {
  before?: Record<string, unknown>;
  after?: Record<string, unknown>;
}

export const ResourceDiff = ({ before, after }: ResourceDiffProps) => {
  return (
    <Flex direction="row" className="gap-4">
      <div className="flex-1">
        <h3 className="text-sm font-medium">Before</h3>
        <pre className="mt-1 rounded-md p-2 text-sm">
          {JSON.stringify(before, null, 2)}
        </pre>
      </div>
      <div className="flex-1">
        <h3 className="text-sm font-medium">After</h3>
        <pre className="mt-1 rounded-md p-2 text-sm">
          {JSON.stringify(after, null, 2)}
        </pre>
      </div>
    </Flex>
  );
};
