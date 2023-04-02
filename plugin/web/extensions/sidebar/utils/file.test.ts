import { expect, test } from "vitest";

import { getExtension, createFileName } from "./file";

test("getExtension", () => {
  expect(getExtension("test.geojson")).toBe("geojson");
  expect(getExtension("test.")).toBe("");
  expect(getExtension("test")).toBe("");
});

test("createFileName", () => {
  expect(createFileName("test", "geojson")).toBe("test.geojson");
  expect(createFileName("", ".czml")).toBe("");
  expect(createFileName("test", "")).toBe("");
  expect(createFileName("", "")).toBe("");
});
