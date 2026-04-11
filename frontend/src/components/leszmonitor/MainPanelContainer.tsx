export const MainPanelContainer = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <main className="flex flex-col gap-6 p-6 w-full h-full">{children}</main>
  );
};
