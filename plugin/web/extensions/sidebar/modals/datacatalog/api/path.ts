export type CatalogItemLike = {
  id: string;
  city?: string;
  pref: string;
  ward?: string;
  group?: string;
  type: string;
  type_en: string;
  type2?: string;
  name?: string;
  /** force to disable making a folder even if type2 is present */
  root?: boolean;
  /** force to make a folder even if type is not special (included in typesWithFolders) */
  root_type?: boolean;
};

export type GroupBy = "city" | "type";

export default function path(
  i: CatalogItemLike,
  customDataset: boolean,
  groupBy: GroupBy,
): string[] {
  return groupBy === "type" ? pathByType(i, customDataset) : pathByCity(i, customDataset);
}

function pathByCity(i: CatalogItemLike, customDataset: boolean): string[] {
  return [
    i.pref,
    ...(i.city ? [i.city] : []),
    ...(i.ward ? [i.ward] : []),
    ...(i.group?.split("/") ?? []),
    ...(!customDataset &&
    (i.root_type ||
      i.type2 ||
      (!i.root && typesWithFolders.includes(i.type_en) && i.pref !== zenkyu))
      ? [i.type || i.type_en]
      : []),
    ...(i.name || i.id).split("/"),
  ];
}

function pathByType(i: CatalogItemLike, customDataset: boolean): string[] {
  return [
    ...(!customDataset ? [i.type || i.type_en] : []),
    i.pref,
    ...((i.ward || i.type2) && i.city ? [i.city] : []),
    ...(i.group?.split("/") ?? []),
    ...(i.name || i.id).split("/"),
  ];
}

// TODO: when root_type is available, these are no longer needed
const zenkyu = "全球データ";
const typesWithFolders = ["usecase", "gen", "fld", "htd", "tnm", "ifld", "urf", "ex"];
