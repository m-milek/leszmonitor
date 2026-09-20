import { Toaster } from "@/components/ui/toast";

export interface ToasterProviderProps {
  children: React.ReactNode;
}

export const ToasterProvider = ({ children }: ToasterProviderProps) => {
  return (
    <Toaster limit={5} timeout={10_000}>
      {children}
    </Toaster>
  );
};
