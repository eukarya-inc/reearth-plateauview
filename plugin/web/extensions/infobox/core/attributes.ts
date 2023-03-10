export default function getAttributes(attributes: any): any {
  if (!attributes || typeof attributes !== "object") return attributes;
  return walk(attributes);
}

function walk(obj: object): object {
  if (!obj || typeof obj !== "object") return obj;

  if (Array.isArray(obj)) {
    return obj.map(o => (typeof o === "object" && o ? walk(o) : o));
  }

  return Object.fromEntries(
    Object.entries(obj)
      .sort((a, b) => a[0].localeCompare(b[0]))
      .map(([k, v]) => {
        const nk = keyMap.get(k);
        const ak = nk ? `${nk}（${k}）` : k;

        if (typeof v === "object" && v) {
          return [ak || k, walk(v)];
        }
        return [ak || k, v];
      }),
  );
}

const keyMap = new Map<string, string>();
