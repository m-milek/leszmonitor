import { Toaster } from "@/components/ui/sonner";

export interface ToasterProviderProps {
  children: React.ReactNode;
}

export const ToasterProvider = ({ children }: ToasterProviderProps) => {
  return (
    <>
      {children}
      <Toaster visibleToasts={5} duration={10_000} />
    </>
  );
};
