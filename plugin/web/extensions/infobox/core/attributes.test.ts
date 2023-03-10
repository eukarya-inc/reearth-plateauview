import { expect, test } from "vitest";

import getAttributes from "./attributes";

test("getAttributes", () => {
  expect(
    getAttributes({
      bbb: {},
      aaa: {
        bbb: "ccc",
        ddd: [{ c: "b" }, { b: "a", a: "" }],
      },
    }),
  ).toEqual({
    aaa: {
      bbb: "ccc",
      ddd: [{ c: "b" }, { a: "", b: "a" }],
    },
    bbb: {},
  });
});
