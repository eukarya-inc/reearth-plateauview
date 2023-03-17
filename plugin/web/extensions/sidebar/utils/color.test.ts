import { expect, test } from "vitest";

import {
  getRGBAFromString,
  generateColorGradient,
  hexToRgb,
  rgbToHex,
  colorToHex,
  colorToRgb,
  rgbToArray,
  setColor,
} from "./color";

test("getRGBAFromString", () => {
  expect(getRGBAFromString("rgba(100, 24, 255, 1)")).toEqual([100, 24, 255, 1]);
  expect(getRGBAFromString("rgba(100,24,255,0.5)")).toEqual([100, 24, 255, 0.5]);
});

test("rgbToArray", () => {
  expect(rgbToArray("rgb(255,0,0)")).toEqual([255, 0, 0]);
  expect(rgbToArray("red")).toEqual([]);
});

test("setColor", () => {
  expect(setColor("rgb(255,0,0)")).toEqual([255, 0, 0]);
  expect(setColor("red")).toEqual([255, 0, 0]);
  expect(setColor("#ff0000")).toEqual([255, 0, 0]);
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
  expect(hexToRgb("#008000")).toEqual([0, 128, 0]);
  expect(hexToRgb("#0000ff")).toEqual([0, 0, 255]);
  expect(hexToRgb("red")).toEqual([0, 0, 0]);
  expect(hexToRgb("something")).toEqual([0, 0, 0]);
});

test("rgbToHex", () => {
  expect(rgbToHex([255, 0, 0])).toEqual("#ff0000");
  expect(rgbToHex([0, 128, 0])).toEqual("#008000");
  expect(rgbToHex([0, 0, 255])).toEqual("#0000ff");
  expect(rgbToHex([])).toEqual("");
  expect(rgbToHex([0, 1])).toEqual("");
});

test("colorToRgb", () => {
  expect(colorToRgb("red")).toEqual([255, 0, 0]);
  expect(colorToRgb("green")).toEqual([0, 128, 0]);
  expect(colorToRgb("blue")).toEqual([0, 0, 255]);
  expect(colorToRgb("something")).toEqual([0, 0, 0]);
  expect(colorToRgb("")).toEqual([0, 0, 0]);
});

test("colorToHex", () => {
  expect(colorToHex("red")).toEqual("#ff0000");
  expect(colorToHex("green")).toEqual("#008000");
  expect(colorToHex("blue")).toEqual("#0000ff");
  expect(colorToHex("something")).toEqual("");
});
