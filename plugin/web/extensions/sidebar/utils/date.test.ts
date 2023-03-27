import { expect, test } from "vitest";

import { formatDateTime } from "./date";

test("formatDateTime", () => {
  expect(formatDateTime("2023-03-25", "12:00:00")).toBe(
    new Date("2023-03-25T12:00:00").toISOString(),
  );
  expect(formatDateTime("xxx", "yyy")).toBe(new Date().toISOString());
});
