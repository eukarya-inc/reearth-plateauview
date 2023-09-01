import { test } from "vitest";

import sortBy, { CatalogItemLike, SortBy } from "./sort";

const base: CatalogItemLike = {
  id: "1",
  pref: "東京都",
  pref_code_i: 1,
  city: "東京都23区",
  ward: "千代田区",
  code: 1,
  type_en: "bldg",
};

const tests: [string, CatalogItemLike, CatalogItemLike, SortBy, boolean][] = [
  [
    "type",
    {
      ...base,
    },
    {
      ...base,
      type_en: "tran",
    },
    "city",
    false,
  ],
  [
    "city with tokyo",
    {
      ...base,
      pref: "北海道",
      city: "札幌市",
      type_en: "usecase",
    },
    {
      ...base,
      pref: "東京都",
      city: "東京都23区",
      type_en: "usecase",
    },
    "city",
    true,
  ],
  [
    "unknown type",
    {
      ...base,
      type_en: "usecase",
    },
    {
      ...base,
      type_en: "hogehoge",
    },
    "city",
    false,
  ],
  [
    "category",
    {
      ...base,
      type_en: "bldg",
      category: "category",
    },
    {
      ...base,
      type_en: "usecase",
    },
    "city",
    true,
  ],
  [
    "two unknown types",
    {
      ...base,
      type_en: "hogehoge",
    },
    {
      ...base,
      type_en: "foobar",
    },
    "city",
    true,
  ],
  [
    "no-group should be upper",
    {
      ...base,
      group: "group",
    },
    {
      ...base,
    },
    "city",
    true,
  ],
];

tests.forEach(([name, a, b, by, reverted]) => {
  test(name, () => {
    const e = expect(sortBy(b, a, by));
    if (reverted) e.toBeLessThan(0);
    else e.toBeGreaterThanOrEqual(0);
  });
});
