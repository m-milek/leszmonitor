import { cn } from "cn";

export interface MainPanelContainerProps {
  className?: string;
  children: React.ReactNode;
}

export const PageContainer = (props: MainPanelContainerProps) => {
  return (
    <div className={cn("flex flex-col gap-4 w-full p-6", props.className)}>
      {props.children}
    </div>
  );
};
