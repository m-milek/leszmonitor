export interface CenterProps {
  children: React.ReactNode;
  centerVertically?: boolean;
  centerHorizontally?: boolean;
}

export function Center({
  children,
  centerVertically = true,
  centerHorizontally = true,
}: Readonly<CenterProps>) {
  return (
    <div
      className={`flex w-full h-full ${centerVertically ? "items-center" : ""} ${centerHorizontally ? "justify-center" : ""}`}
    >
      {children}
    </div>
  );
}
