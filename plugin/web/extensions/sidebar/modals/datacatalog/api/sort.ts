export type CatalogItemLike = {
  id: string;
  pref?: string;
  city?: string;
  ward?: string;
  type_en: string;
  code: number;
  group?: string;
  category?: string;
  // internal
  pref_code_i: number;
  order?: number;
};

// tag is not used for sorting, but it is refered from data catalog components
export type SortBy = "city" | "type" | "tag";

export default function sortBy(a: CatalogItemLike, b: CatalogItemLike, sort: SortBy): number {
  return sort === "type"
    ? sortByType(a, b) || sortByCity(a, b) || sortByOrder(a.order, b.order)
    : sortByCity(a, b) || sortByType(a, b) || sortByOrder(a.order, b.order);
}

function sortByCity(a: CatalogItemLike, b: CatalogItemLike): number {
  return clamp(
    (a.pref === zenkyu ? 0 : 1) - (b.pref === zenkyu ? 0 : 1) || // items whose prefecture is zenkyu is upper
      (a.pref === tokyo ? 0 : 1) - (b.pref === tokyo ? 0 : 1) || // items whose prefecture is tokyo is upper
      a.pref_code_i - b.pref_code_i ||
      (a.ward ? 0 : 1) - (b.ward ? 0 : 1) || // items that have a ward is upper
      (a.city ? 0 : 1) - (b.city ? 0 : 1) || // items that have a city is upper
      (b.group ? 0 : 1) - (a.group ? 0 : 1) || // items that have no groups is upper
      a.code - b.code,
  );
}

function sortByType(a: CatalogItemLike, b: CatalogItemLike): number {
  const ai = a.category ? -1 : types.indexOf(a.type_en);
  const bi = b.category ? -1 : types.indexOf(b.type_en);

  if (ai === -1 && bi !== -1) return 1;
  if (ai !== -1 && bi === -1) return -1;
  if (ai === -1 && bi === -1) {
    return (a.category || a.type_en).localeCompare(b.category || b.type_en);
  }

  return clamp(ai - bi);
}

function sortByOrder(a: number | undefined, b: number | undefined): number {
  return clamp(Math.min(0, a ?? 0) - Math.min(0, b ?? 0));
}

function clamp(n: number): number {
  return Math.max(-1, Math.min(1, n));
}

const zenkyu = "全球データ";
const tokyo = "東京都";

const types = [
  "bldg",
  "tran",
  "brid",
  "rail",
  "veg",
  "frn",
  "luse",
  "lsld",
  "urf",
  "fld",
  "tnm",
  "htd",
  "ifld",
  "gen",
  "shelter",
  "landmark",
  "station",
  "emergency_route",
  "railway",
  "park",
  "border",
  "usecase",
];
