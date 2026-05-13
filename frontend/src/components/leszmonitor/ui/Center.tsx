export interface CenterProps {
  children: React.ReactNode;
}

export function Center({ children }: Readonly<CenterProps>) {
  return (
    <div className="flex items-center justify-center w-full h-full">
      {children}
    </div>
  );
}
