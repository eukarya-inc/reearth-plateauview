import { expect, test } from "vitest";

import { getExtension, joinNameWithExtension } from "./file";

test("getExtension", () => {
  expect(getExtension("test.geojson")).toBe("geojson");
  expect(getExtension("test.")).toBe("");
  expect(getExtension("test")).toBe("");
});

test("joinNameWithExtension", () => {
  expect(joinNameWithExtension("test", "geojson")).toBe("test.geojson");
  expect(joinNameWithExtension("", ".czml")).toBe("");
  expect(joinNameWithExtension("test", "")).toBe("");
  expect(joinNameWithExtension("", "")).toBe("");
});
