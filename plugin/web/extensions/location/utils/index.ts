import { PostMessageProps } from "../types";

export function postMsg({ action, payload }: PostMessageProps) {
  parent.postMessage(
    {
      action,
      payload,
    },
    "*",
  );
}
export const distances = [
  1, 2, 3, 5, 10, 20, 30, 50, 100, 200, 300, 500, 1000, 2000, 3000, 5000, 10000, 20000, 30000,
  50000, 100000, 200000, 300000, 500000, 1000000, 2000000, 3000000, 5000000, 10000000, 20000000,
  30000000, 50000000,
];
