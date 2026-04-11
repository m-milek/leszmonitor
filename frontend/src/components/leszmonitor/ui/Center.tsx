export interface CenterProps {
  children: React.ReactNode;
}

export function Center({ children }: CenterProps) {
  return (
    <div className="flex items-center justify-center w-full h-full">
      {children}
    </div>
  );
}
