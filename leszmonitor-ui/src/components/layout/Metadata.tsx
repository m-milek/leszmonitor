import type { MetadataResponse } from "@/features/instance/metadata-api";
import { Center } from "@/components/common/Center";
import { LucideInfo } from "lucide-react";
import { Flex } from "@/components/common/Flex";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

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
    <Center className="p-2">
      <Flex className="gap-2 items-center">
        <span className="text-sm">{data?.version}</span>
        <Popover>
          <PopoverTrigger>
            <LucideInfo className="size-4 cursor-pointer" />
          </PopoverTrigger>
          <PopoverContent className="w-80">
            <MetadataTable data={data} />
          </PopoverContent>
        </Popover>
      </Flex>
    </Center>
  );
};
