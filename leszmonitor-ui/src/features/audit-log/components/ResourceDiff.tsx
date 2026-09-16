export interface ResourceDiffProps {
  before?: string;
  after?: string;
}

const parseRecursively = (obj: any): any => {
  if (typeof obj === "string") {
    try {
      const parsed = JSON.parse(obj);
      if (parsed !== null && typeof parsed === "object") {
        return parseRecursively(parsed);
      }
      return obj;
    } catch {
      return obj;
    }
  } else if (Array.isArray(obj)) {
    return obj.map(parseRecursively);
  } else if (obj !== null && typeof obj === "object") {
    const newObj: any = {};
    for (const key in obj) {
      newObj[key] = parseRecursively(obj[key]);
    }
    return newObj;
  }
  return obj;
};

const safeParseJSON = (jsonString?: string) => {
  if (!jsonString) return "—";
  try {
    let parsed = JSON.parse(jsonString);
    if (parsed === null) return "—";
    parsed = parseRecursively(parsed);
    return JSON.stringify(parsed, null, 2);
  } catch (e) {
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
