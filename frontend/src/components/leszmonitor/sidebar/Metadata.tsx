import { Center } from "@/components/leszmonitor/ui/Center.tsx";
import type { MetadataResponse } from "@/lib/data/metadata-api.ts";
import { LucideInfo } from "lucide-react";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover.tsx";

interface MetadataProps {
  data?: MetadataResponse;
}

const MetadataTable = ({ data }: { data?: MetadataResponse }) => (
  <table className="w-full">
    <tbody>
      <tr>
        <td className="font-medium">Version</td>
        <td>{data?.version}</td>
      </tr>
      <tr>
        <td className="font-medium">CI Build Number</td>
        <td>{data?.ciBuildNumber}</td>
      </tr>
      <tr>
        <td className="font-medium">Git Commit</td>
        <td>{data?.gitCommit}</td>
      </tr>
      <tr>
        <td className="font-medium">Image Tag</td>
        <td>{data?.imageTag}</td>
      </tr>
    </tbody>
  </table>
);

export const Metadata = ({ data }: MetadataProps) => {
  return (
    <Center className="p-1">
      <Flex className="gap-2">
        {data?.version}
        <Popover>
          <PopoverTrigger>
            <LucideInfo className="cursor-pointer" />
          </PopoverTrigger>
          <PopoverContent className="w-80">
            <MetadataTable data={data} />
          </PopoverContent>
        </Popover>
      </Flex>
    </Center>
  );
};
