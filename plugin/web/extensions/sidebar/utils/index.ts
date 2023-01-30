import { PostMessageProps } from "@web/extensions/sidebar/types";
import { cloneDeep, mergeWith } from "lodash";

export function postMsg({ action, payload }: PostMessageProps) {
  parent.postMessage(
    {
      action,
      payload,
    },
    "*",
  );
}

export function mergeProperty(a: any, b: any) {
  const a2 = cloneDeep(a);
  return mergeWith(
    a2,
    b,
    (s: any, v: any, _k: string | number | symbol, _obj: any, _src: any, stack: { size: number }) =>
      stack.size > 0 || Array.isArray(v) ? v ?? s : undefined,
  );
}

export function rgbToHex(r: number, g: number, b: number) {
  return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

export function hexToRgb(hex: string) {
  const result: RegExpExecArray | null = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result
    ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16),
      }
    : { r: 0, g: 0, b: 0 };
}

// returns an array of startColor, colors between according to steps, and endColor
export function getColorsGradients(startColor: string, endColor: string, steps: number) {
  const Gradients = [];

  Gradients.push(startColor);

  const startColorRgb = hexToRgb(startColor);
  const endColorRgb = hexToRgb(endColor) ?? 0;

  const rInc = Math.round((endColorRgb.r - startColorRgb.r) / steps + 1);
  const gInc = Math.round((endColorRgb.g - startColorRgb.g) / steps + 1);
  const bInc = Math.round((endColorRgb.b - startColorRgb.b) / steps + 1);

  for (let i = 0; i < steps - 2; i++) {
    startColorRgb.r += rInc;
    startColorRgb.g += gInc;
    startColorRgb.b += bInc;

    Gradients.push(rgbToHex(startColorRgb.r, startColorRgb.g, startColorRgb.b));
  }
  Gradients.push(endColor);

  return Gradients;
}
