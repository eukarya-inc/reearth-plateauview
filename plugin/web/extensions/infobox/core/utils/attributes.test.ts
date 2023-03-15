import { expect, test, vi } from "vitest";

import { attributesMap, getAttributes, getFields } from "./attributes";
import type { Json } from "./json";

test("getAttributes", () => {
  const src: Json = {
    bbb: {},
    aaa: {
      bbb: "ccc",
      ddd: [{ c: "b" }, { b: "a", a: "" }],
    },
  };
  expect(flatKeys(src)).toEqual([
    "",
    "bbb",
    "aaa",
    "aaa.bbb",
    "aaa.ddd",
    "aaa.ddd.0",
    "aaa.ddd.0.c",
    "aaa.ddd.1",
    "aaa.ddd.1.b",
    "aaa.ddd.1.a",
  ]);

  const actual = getAttributes(src);
  expect(flatKeys(actual)).toEqual([
    "",
    "AAA（aaa）",
    "AAA（aaa）.bbb",
    "AAA（aaa）.DDD（ddd）",
    "AAA（aaa）.DDD（ddd）.0",
    "AAA（aaa）.DDD（ddd）.0.c",
    "AAA（aaa）.DDD（ddd）.1",
    "AAA（aaa）.DDD（ddd）.1.a",
    "AAA（aaa）.DDD（ddd）.1.b",
    "bbb",
  ]);
});

test("getFields", () => {
  expect(
    getFields(
      { a: { b: "aaa" } },
      {
        A: ["a", "b"],
        B: { keys: ["a", "b"], map: (v: string) => v + "!" },
        C: ["b"],
      },
    ),
  ).toEqual({
    A: "aaa",
    B: "aaa!",
  });
});

test("attributesMap", () => {
  expect(attributesMap.get("ddd")).toBe("DDD");
});

function flatKeys(obj: Json, parentKey?: string): string[] {
  if (typeof obj !== "object" || !obj) return [parentKey || ""];
  return [
    parentKey || "",
    ...Object.entries(obj).flatMap(([k, v]) =>
      flatKeys(v, `${parentKey ? `${parentKey}.` : ""}${k}`),
    ),
  ];
}

vi.mock("./attributes.csv?raw", () => ({
  default: "ddd,DDD\naaa,AAA\n",
}));
