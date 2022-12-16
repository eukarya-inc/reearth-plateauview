import { cloneDeep, mergeWith } from "lodash";

import { PostMessageActionType } from "../types";

export function postMsg(action: PostMessageActionType, payload?: any) {
  parent.postMessage({
    action,
    payload,
  });
}

export function mergeProperty(a: any, b: any) {
  const a2 = cloneDeep(a);
  return mergeWith(
    a2,
    b,
    (s: any, v: any, _k: string | number | symbol, _obj: any, _src: any, stack: { size: number }) =>
      stack.size > 0 || Array.isArray(v) ? v ?? s : undefined,
  );
}
