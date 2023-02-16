import { omit, makeTree, mapTree } from "./utils";

export type RawDataCatalogItem = {
  id: string;
  name: string;
  pref: string;
  city: string;
  city_en: string;
  city_code: string;
  ward?: string;
  ward_en?: string;
  ward_code?: string;
  type: string;
  type_en: string;
  format: string;
  layers?: string[];
  url: string;
  desc: string;
  year: number;
  config?: any;
  // bldg only fields
  bldg_low_texture_url?: string;
  bldg_no_texture_url?: string;
  search_index?: string;
};

export type RawDataCatalogGroup = {
  name: string;
  children: (RawDataCatalogItem | RawDataCatalogGroup)[];
};

export type GroupBy = "city" | "type" | "tag"; // Tag not implemented yet

export async function getDataCatalog(base: string): Promise<RawDataCatalogItem[]> {
  const res = await fetch(base + "/datacatalog");
  if (res.status !== 200) {
    throw new Error("failed to fetch data catalog");
  }
  return res.json();
}

export function getDataCatalogTree(
  items: RawDataCatalogItem[],
  groupBy: GroupBy,
  q?: string | undefined,
): (RawDataCatalogGroup | RawDataCatalogItem)[] {
  const allItems = filter(q, items)
    .map(i => ({
      ...i,
      pref: i.pref === "全国" || i.pref === "全球" ? "全球データ" : i.pref,
      path: path(i, groupBy),
      code: i.ward_code ? parseInt(i.ward_code) : parseInt(i.city_code),
    }))
    .sort((a, b) => sortBy(a, b, groupBy));

  return mapTree(makeTree(allItems), (item): RawDataCatalogGroup | RawDataCatalogItem =>
    item.item
      ? omit(item.item, "path", "code")
      : {
          name: item.name,
          children: [],
        },
  );
}

function path(i: RawDataCatalogItem, groupBy: GroupBy): string[] {
  return groupBy === "type"
    ? [i.type, i.pref, ...(i.ward ? [i.city] : []), ...i.name.split("/")]
    : [i.pref, i.city, ...(i.ward ? [i.ward] : []), ...i.name.split("/")];
}

function sortBy(
  a: RawDataCatalogItem & { code: number },
  b: RawDataCatalogItem & { code: number },
  sort: GroupBy,
): number {
  return sort === "type"
    ? sortByType(a, b) || sortByCity(a, b)
    : sortByCity(a, b) || sortByType(a, b);
}

function sortByCity(
  a: RawDataCatalogItem & { code: number },
  b: RawDataCatalogItem & { code: number },
): number {
  return (
    (zenkoku.includes(b.pref) ? 1 : 0) - (zenkoku.includes(a.pref) ? 1 : 0) ||
    (b.pref === tokyo ? 1 : 0) - (a.pref === tokyo ? 1 : 0) ||
    (!a.city ? 1 : 0) - (!b.city ? 1 : 0) ||
    (!a.ward ? 1 : 0) - (!b.ward ? 1 : 0) ||
    a.code - b.code ||
    types.indexOf(a.type_en) - types.indexOf(b.type_en)
  );
}
function sortByType(a: RawDataCatalogItem, b: RawDataCatalogItem): number {
  return types.indexOf(a.type_en) - types.indexOf(b.type_en);
}

function filter(q: string | undefined, items: RawDataCatalogItem[]): RawDataCatalogItem[] {
  if (!q) return items;
  return items.filter(
    i => i.name.includes(q) || i.pref.includes(q) || i.city?.includes(q) || i.ward?.includes(q),
  );
}

const zenkoku = ["全国", "全球", "全球データ"];
const tokyo = "東京都";
const types = [
  "bldg",
  "tran",
  "veg",
  "frn",
  "luse",
  "lsld",
  "urf",
  "fld",
  "tnm",
  "htd",
  "ifld",
  "shelter",
  "landmark",
  "station",
  "emergency_route",
  "railway",
  "park",
  "border",
  "usecase",
];
