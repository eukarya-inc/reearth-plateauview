import { expect, test } from "vitest";

import { formatDateTime } from "./date";

test("formatDateTime", () => {
  expect(formatDateTime("2023-03-25", "12:00:00")).toBe("2023-03-25T03:00:00.000Z");
  expect(formatDateTime("xxx", "yyy")).toBe(new Date().toISOString());
});
