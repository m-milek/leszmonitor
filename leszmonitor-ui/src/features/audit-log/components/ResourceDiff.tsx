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

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
      <div className="flex flex-col min-w-0">
        <h3 className="text-sm font-medium mb-1">Before</h3>
        <pre className="p-4 bg-muted/50 rounded-lg text-sm overflow-auto flex-1 border border-border">
          {beforePrettyJSON}
        </pre>
      </div>
      <div className="flex flex-col min-w-0">
        <h3 className="text-sm font-medium mb-1">After</h3>
        <pre className="p-4 bg-muted/50 rounded-lg text-sm overflow-auto flex-1 border border-border">
          {afterPrettyJSON}
        </pre>
      </div>
    </div>
  );
};
