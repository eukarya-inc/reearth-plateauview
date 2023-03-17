import tinycolor from "tinycolor2";

export type RGBA = [r: number, g: number, b: number, a: number];

export const getRGBAFromString = (rgbaStr: string | undefined): RGBA | undefined => {
  const matches = rgbaStr?.match(/rgba\((\d*), *(\d*), *(\d*), *((\d|\.)*)\)/)?.slice(0, -1);
  return matches ? (matches.slice(1).map(m => Number(m)) as RGBA) : undefined;
};

export const rgbaToString = (rgba: RGBA) => `rgba(${rgba.join(",")})`;

export const RGB_COLOR_REGEX = /^rgb\((\d{1,3}),\s*(\d{1,3}),\s*(\d{1,3})\)$/;
export const HEX_COLOR_REGEX = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/;

export const rgbToArray = (rgb: string): number[] => {
  return RGB_COLOR_REGEX.test(rgb)
    ? rgb
        .replace(/^rgb\(|\s+|\)$/g, "")
        .split(",")
        .map(Number)
    : [];
};

export const setColor = (color: string) => {
  return HEX_COLOR_REGEX.test(color)
    ? hexToRgb(color)
    : RGB_COLOR_REGEX.test(color)
    ? rgbToArray(color)
    : colorToRgb(color);
};

export const generateColorGradient = (
  colorStart: string,
  colorEnd: string,
  colorCount: number,
): string[] => {
  const startRGB = setColor(colorStart);
  const endRGB = setColor(colorEnd);

  const stepR = (endRGB[0] - startRGB[0]) / (colorCount - 1);
  const stepG = (endRGB[1] - startRGB[1]) / (colorCount - 1);
  const stepB = (endRGB[2] - startRGB[2]) / (colorCount - 1);

  const gradientColors = [];
  for (let i = 0; i < colorCount; i++) {
    const r = Math.round(startRGB[0] + stepR * i);
    const g = Math.round(startRGB[1] + stepG * i);
    const b = Math.round(startRGB[2] + stepB * i);
    const color = rgbToHex([r, g, b]);
    gradientColors.push(color);
  }

  return gradientColors;
};

export const hexToRgb = (hex: string): number[] => {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return !isNaN(r) && !isNaN(g) && !isNaN(b) ? [r, g, b] : [0, 0, 0];
};

export const rgbToHex = (rgb: number[]): string => {
  if (rgb.length < 3) return "";
  const r = rgb[0].toString(16).padStart(2, "0");
  const g = rgb[1].toString(16).padStart(2, "0");
  const b = rgb[2].toString(16).padStart(2, "0");
  return `#${r}${g}${b}`;
};

export const colorToRgb = (color: string): number[] => hexToRgb(colorToHex(color));

export const colorToHex = (color: string): string => {
  return tinycolor(color).toHexString();
};
