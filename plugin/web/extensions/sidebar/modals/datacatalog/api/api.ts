import type {
  DataCatalogGroup,
  DataCatalogItem,
  DataCatalogTreeItem,
} from "@web/extensions/sidebar/core/types";
import { omit } from "lodash-es";

import path, { GroupBy } from "./path";
import sortBy from "./sort";
import { makeTree, mapTree } from "./utils";

// TODO: REFACTOR: CONFUSING REEXPORT
export type { DataCatalogItem, DataCatalogGroup, DataCatalogTreeItem };

export type RawDataCatalogTreeItem = RawDataCatalogGroup | RawDataCatalogItem;

export type { GroupBy } from "./path";

export type RawDataCatalogGroup = {
  id: string;
  name: string;
  desc?: string;
  children: RawDataCatalogTreeItem[];
};

type RawRawDataCatalogItem = {
  id: string;
  itemId?: string;
  name?: string;
  pref: string;
  pref_code?: string;
  pref_code_i: number;
  city?: string;
  city_en?: string;
  city_code?: string;
  city_code_i: number;
  ward?: string;
  ward_en?: string;
  ward_code?: string;
  ward_code_i: number;
  type: string;
  type_en: string;
  type2?: string;
  type2_en?: string;
  format: string;
  layers?: string[] | string;
  layer?: string[] | string;
  url: string;
  desc: string;
  year: number;
  tags?: { type: "type" | "location"; value: string }[];
  openDataUrl?: string;
  config?: {
    data?: {
      name: string;
      type: string;
      url: string;
      layers?: string[] | string;
      layer?: string[] | string;
    }[];
  };
  order?: number;
  group?: string;
  /** force to disable making a folder even if type2 is present */
  root?: boolean;
  /** force to make a folder even if type is not special (included in typesWithFolders) */
  root_type?: boolean;
  /** alias of type that is used as a folder name */
  category?: string;
  infobox?: boolean;

  // bldg only fields
  search_index?: string;

  // internal
  path?: string[];
  code: number;
};

export type DataSource = "plateau" | "custom";

export type RawDataCatalogItem = Omit<RawRawDataCatalogItem, "layers" | "layer" | "config"> & {
  layers?: string[];
  config?: {
    data?: {
      name: string;
      type: string;
      url: string;
      layer?: string[];
    }[];
  };
};

export async function getDataCatalog(
  base: string,
  project?: string,
  dataSource?: DataSource,
): Promise<RawDataCatalogItem[]> {
  const res = await fetch(`${base}/datacatalog${project ? `/${project}` : ""}`);
  if (!res.ok) {
    throw new Error("failed to fetch data catalog");
  }

  const data: RawRawDataCatalogItem[] = await res.json();
  return data.map(d => modifyDataCatalog(d, dataSource));
}

export function modifyDataCatalog(
  d: Omit<RawRawDataCatalogItem, "pref_code_i" | "city_code_i" | "ward_code_i" | "tags" | "code">,
  dataSource?: DataSource,
): RawDataCatalogItem & { dataSource?: DataSource } {
  const pref = d.pref === "全国" || d.pref === "全球" ? zenkyu : d.pref;
  const pref_code = d.pref === "全国" || d.pref === "全球" || d.pref === zenkyu ? "0" : d.pref_code;
  const pref_code_i = parseInt(pref_code ?? "");
  const city_code_i = parseInt(d.city_code ?? "");
  const ward_code_i = parseInt(d.ward_code ?? "");
  return {
    ...omit(d, ["layers", "layer", "config"]),
    pref,
    pref_code,
    pref_code_i,
    city_code_i,
    ward_code_i,
    code: !isNaN(ward_code_i)
      ? ward_code_i
      : !isNaN(city_code_i)
      ? city_code_i
      : !isNaN(pref_code_i)
      ? pref_code_i * 1000
      : pref === zenkyu
      ? 0
      : 99999,
    tags: [
      { type: "type", value: d.type },
      ...(d.type2 ? [{ type: "type", value: d.type2 } as const] : []),
      ...(d.city ? [{ type: "location", value: d.city } as const] : []),
      ...(d.ward ? [{ type: "location", value: d.ward } as const] : []),
    ],
    ...(d.layers || d.layer ? { layers: [...getLayers(d.layers), ...getLayers(d.layer)] } : {}),
    ...(d.config
      ? {
          config: {
            ...(d.config.data
              ? {
                  data: d.config.data.map(dd => ({
                    ...omit(dd, ["layers", "layer"]),
                    layer: [...getLayers(dd.layers), ...getLayers(dd.layer)],
                  })),
                }
              : {}),
          },
        }
      : {}),
    dataSource,
  };
}

// TODO: REFACTOR: confusing typing
export function getDataCatalogTree(
  items: DataCatalogItem[],
  groupBy: GroupBy,
  customDataset: boolean,
  q?: string | undefined,
): DataCatalogTreeItem[] {
  return getRawDataCatalogTree(items, groupBy, customDataset, q) as DataCatalogTreeItem[];
}

export function getRawDataCatalogTree(
  items: RawDataCatalogItem[],
  groupBy: GroupBy,
  customDataset: boolean,
  q?: string | undefined,
): (RawDataCatalogGroup | RawDataCatalogItem)[] {
  return mapTree(
    makeTree(sortInternal(items, groupBy, customDataset, q)),
    (item): RawDataCatalogGroup | RawDataCatalogItem =>
      item.item ?? {
        id: item.id,
        name: item.name,
        desc: item.desc,
        children: [],
      },
  );
}

type InternalDataCatalogItem = RawDataCatalogItem & {
  path: string[];
};

function sortInternal(
  items: RawDataCatalogItem[],
  groupBy: GroupBy,
  customDataset: boolean,
  q?: string | undefined,
): InternalDataCatalogItem[] {
  return filter(q, items)
    .map(
      (i): InternalDataCatalogItem => ({
        ...i,
        path: path(i, customDataset, groupBy),
      }),
    )
    .sort((a, b) => sortBy(a, b, groupBy));
}

function filter(q: string | undefined, items: RawDataCatalogItem[]): RawDataCatalogItem[] {
  if (!q) return items;
  return items.filter(
    i =>
      i.name?.includes(q) ||
      i.pref.includes(q) ||
      i.city?.includes(q) ||
      i.ward?.includes(q) ||
      i.type_en === "folder",
  );
}

function getLayers(layers?: string[] | string): string[] {
  return layers ? (typeof layers === "string" ? layers.split(/, */).filter(Boolean) : layers) : [];
}

const zenkyu = "全球データ";
