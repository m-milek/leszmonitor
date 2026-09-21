import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Field, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useTheme } from "next-themes";

const themes = [
  { value: "light", label: "Light" },
  { value: "dark", label: "Dark" },
  { value: "system", label: "System" },
];

export function UserSettingsPage() {
  const { theme, setTheme } = useTheme();

  return (
    <PageContainer>
      <TypographyH1>Settings</TypographyH1>
      <Field className="max-w-xs">
        <FieldLabel>Theme</FieldLabel>
        <Select
          value={theme ?? "system"}
          onValueChange={(value) => setTheme(value ?? "system")}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent alignItemWithTrigger={false}>
            {themes.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </Field>
    </PageContainer>
  );
}
