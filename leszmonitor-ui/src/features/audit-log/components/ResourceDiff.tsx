import ReactDiffViewer from "react-diff-viewer-continued";
import { useTheme } from "next-themes";

export interface ResourceDiffProps {
  before?: string;
  after?: string;
}

const parseRecursively = (value: unknown): unknown => {
  if (typeof value === "string") {
    try {
      const parsed: unknown = JSON.parse(value);
      if (parsed !== null && typeof parsed === "object") {
        return parseRecursively(parsed);
      }
      return value;
    } catch {
      return value;
    }
  }

  if (Array.isArray(value)) {
    return value.map(parseRecursively);
  }

  if (value !== null && typeof value === "object") {
    const result: Record<string, unknown> = {};
    for (const [key, entry] of Object.entries(value)) {
      result[key] = parseRecursively(entry);
    }
    return result;
  }

  return value;
};

const safeParseJSON = (jsonString?: string) => {
  if (!jsonString) return "—";
  try {
    const parsed: unknown = JSON.parse(jsonString);
    if (parsed === null) return "—";
    return JSON.stringify(parseRecursively(parsed), null, 2);
  } catch {
    return jsonString;
  }
};

export const ResourceDiff = ({ before, after }: ResourceDiffProps) => {
  const beforePrettyJSON = safeParseJSON(before);
  const afterPrettyJSON = safeParseJSON(after);

  const { resolvedTheme } = useTheme();

  return (
    <ReactDiffViewer
      hideSummary
      useDarkTheme={resolvedTheme === "dark"}
      oldValue={beforePrettyJSON}
      newValue={afterPrettyJSON}
      splitView={true}
    />
  );
};
