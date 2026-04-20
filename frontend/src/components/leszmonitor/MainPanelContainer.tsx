import { cn } from "@/lib/utils.ts";

export interface MainPanelContainerProps {
  className?: string;
  children: React.ReactNode;
}

export const MainPanelContainer = (props: MainPanelContainerProps) => {
  return (
    <main
      className={cn("flex flex-col gap-6 p-6 w-full h-full", props.className)}
    >
      {props.children}
    </main>
  );
};
