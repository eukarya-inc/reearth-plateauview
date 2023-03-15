import { get } from "lodash";

import attributesData from "./attributes.csv?raw";
import type { Json, JsonArray, JsonObject } from "./json";

export const attributesMap = new Map<string, string>();

attributesData
  .split("\n")
  .map(l => l.split(","))
  .forEach(l => {
    if (!l || !l[0] || !l[1] || typeof l[0] !== "string" || typeof l[1] !== "string") return;
    attributesMap.set(l[0], l[1]);
  });

export function getAttributes(attributes: Json): Json {
  if (!attributes || typeof attributes !== "object") return attributes;
  return walk(attributes, attributesMap);
}

function walk(obj: JsonObject | JsonArray, keyMap?: Map<string, string>): JsonObject | JsonArray {
  if (!obj || typeof obj !== "object") return obj;

  if (Array.isArray(obj)) {
    return obj.map(o => (typeof o === "object" && o ? walk(o) : o));
  }

  return Object.fromEntries(
    Object.entries(obj)
      .sort((a, b) => a[0].localeCompare(b[0]))
      .map(([k, v]) => {
        const nk = keyMap?.get(k);
        const ak = nk ? `${nk}（${k}）` : k;

        if (typeof v === "object" && v) {
          return [ak || k, walk(v, keyMap)];
        }
        return [ak || k, v];
      }),
  );
}

export function getFields(
  attributes: Json,
  keys: Record<string, string[] | { keys: string[]; map: (v: any) => any }>,
): Record<string, any> {
  return Object.fromEntries(
    Object.entries(keys)
      .map(([k, v]) => {
        const keys = Array.isArray(v) ? v : v.keys;
        const map = Array.isArray(v) ? undefined : v.map;
        const nv = get(attributes, keys);
        return nv ? [k, map?.(nv) ?? nv] : null;
      })
      .filter((e): e is [string, any] => !!e),
  );
}
