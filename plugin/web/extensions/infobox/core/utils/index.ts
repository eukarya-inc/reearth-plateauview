import { ActionType } from "../../types";

export function postMsg(action: ActionType, payload?: any) {
  if (parent === window) return;
  parent.postMessage({
    action,
    payload,
  });
}

export function generateId() {
  return "xxxxxxxxxxxxxxxxxxxxxxxxxx".replace(/[x]/g, function () {
    return ((Math.random() * 16) | 0).toString(16);
  });
}
