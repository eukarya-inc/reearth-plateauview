import { test } from "vitest";

import path, { CatalogItemLike, GroupBy } from "./path";

const tests: [CatalogItemLike, boolean, GroupBy, string[]][] = [
  // city
  [
    // regular
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "type",
      type_en: "type_en",
      name: "city_regular",
    },
    false,
    "city",
    ["pref", "city", "city_regular"],
  ],
  [
    // splitted name
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "type",
      type_en: "type_en",
      name: "a/city_splitted_name",
    },
    false,
    "city",
    ["pref", "city", "a", "city_splitted_name"],
  ],
  [
    // ward
    {
      id: "id",
      pref: "pref",
      city: "city",
      ward: "ward",
      type: "type",
      type_en: "type_en",
      name: "city_ward",
    },
    false,
    "city",
    ["pref", "city", "ward", "city_ward"],
  ],
  [
    // group
    {
      id: "id",
      pref: "pref",
      city: "city",
      group: "a/b/c",
      type: "タイプ",
      type_en: "type",
      name: "city_group",
    },
    false,
    "city",
    ["pref", "city", "a", "b", "c", "city_group"],
  ],
  [
    // root
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "gen!",
      type_en: "gen", // gen is included in typesWithFolders
      root: true,
      name: "city_root",
    },
    false,
    "city",
    ["pref", "city", "city_root"],
  ],
  [
    // root_type
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "タイプ",
      type_en: "type",
      root_type: true,
      name: "city_root_type",
    },
    false,
    "city",
    ["pref", "city", "タイプ", "city_root_type"],
  ],
  [
    // usecase
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "ユースケース",
      type_en: "usecase",
      name: "city_usecase",
      root_type: true, // normal usecase data has root_type
    },
    false,
    "city",
    ["pref", "city", "ユースケース", "city_usecase"],
  ],
  [
    // zenkyu
    {
      id: "id",
      pref: "全球データ",
      type: "タイプ",
      type_en: "type",
      name: "city_zenkyu",
    },
    false,
    "city",
    ["全球データ", "city_zenkyu"],
  ],
  // type
  [
    // regular
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "type",
      type_en: "type_en",
      name: "type_regular",
    },
    false,
    "type",
    ["type", "pref", "type_regular"],
  ],
  [
    // splitted name
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "type",
      type_en: "type_en",
      name: "a/type_splitted_name",
    },
    false,
    "type",
    ["type", "pref", "a", "type_splitted_name"],
  ],
  [
    // ward
    {
      id: "id",
      pref: "pref",
      city: "city",
      ward: "ward",
      type: "type",
      type_en: "type_en",
      name: "city_ward",
    },
    false,
    "type",
    ["type", "pref", "city", "city_ward"],
  ],
  [
    // group
    {
      id: "id",
      pref: "pref",
      city: "city",
      group: "a/b/c",
      type: "タイプ",
      type_en: "type",
      name: "type_group",
    },
    false,
    "type",
    ["タイプ", "pref", "a", "b", "c", "type_group"],
  ],
  [
    // usecase
    {
      id: "id",
      pref: "pref",
      city: "city",
      type: "ユースケース",
      type_en: "usecase",
      name: "type_usecase",
    },
    false,
    "type",
    ["ユースケース", "pref", "type_usecase"],
  ],
  [
    // zenkyu
    {
      id: "id",
      pref: "全球データ",
      type: "タイプ",
      type_en: "type",
      name: "type_zenkyu",
    },
    false,
    "type",
    ["タイプ", "全球データ", "type_zenkyu"],
  ],
];

test("path", () => {
  tests.forEach(([i, customDataset, groupBy, expected]) => {
    expect(path(i, customDataset, groupBy)).toEqual(expected);
  });
});
