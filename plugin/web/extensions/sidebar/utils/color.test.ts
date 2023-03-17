import { expect, test } from "vitest";

import { getRGBAFromString, generateColorGradient, hexToRgb, rgbToHex } from "./color";

test("getRGBAFromString", () => {
  expect(getRGBAFromString("rgba(100, 24, 255, 1)")).toEqual([100, 24, 255, 1]);
  expect(getRGBAFromString("rgba(100,24,255,0.5)")).toEqual([100, 24, 255, 0.5]);
});

test("generateColorGradient", () => {
  expect(generateColorGradient("#ff0000", "#0000ff", 5)).toEqual([
    "#ff0000",
    "#bf0040",
    "#800080",
    "#4000bf",
    "#0000ff",
  ]);
});

test("hexToRgb", () => {
  expect(hexToRgb("#ff0000")).toEqual([255, 0, 0]);
  expect(hexToRgb("#00ff00")).toEqual([0, 255, 0]);
  expect(hexToRgb("#0000ff")).toEqual([0, 0, 255]);
});

test("rgbToHex", () => {
  expect(rgbToHex([255, 0, 0])).toEqual("#ff0000");
  expect(rgbToHex([0, 255, 0])).toEqual("#00ff00");
  expect(rgbToHex([0, 0, 255])).toEqual("#0000ff");
});
