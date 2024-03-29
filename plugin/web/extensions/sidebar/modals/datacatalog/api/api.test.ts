import { readFileSync, writeFileSync } from "fs";

import { expect, test } from "vitest";

import {
  // getDataCatalog,
  getRawDataCatalogTree,
  modifyDataCatalog,
  type RawDataCatalogItem,
  type RawDataCatalogTreeItem,
} from "./api";

test.skip("getRawDataCatalogTree", async () => {
  // const r = await getDataCatalog("https://", "");
  const r = JSON.parse(readFileSync("datacatalog.json", "utf8"));
  const d = getRawDataCatalogTree(r, "city", false, "");
  const s = print(d);
  writeFileSync("test.txt", s);
});

function print(c: RawDataCatalogTreeItem[], depth?: number): string {
  let s = "";

  c.forEach(c => {
    s += " ".repeat(depth || 0) + " " + c.name + "\n";
    if ("children" in c) {
      s += print(c.children, (depth || 0) + 1);
    }
  });

  return s;
}

test("modifyDataCatalog", () => {
  const d = {
    id: "a",
    name: "name",
    pref: "pref",
    desc: "",
    format: "",
    url: "",
    type: "ユースケース",
    type_en: "usecase",
    type2: "type2",
    type2_en: "type2",
    city: "city",
    city_code: "11111",
    ward: "ward",
    ward_code: "11112",
    city_en: "city",
    year: 2022,
    root_type: true,
  };
  expect(modifyDataCatalog(d)).toEqual({
    ...d,
    pref_code_i: NaN,
    city_code_i: 11111,
    ward_code_i: 11112,
    code: 11112,
    tags: [
      { type: "type", value: "ユースケース" },
      { type: "type", value: "type2" },
      { type: "location", value: "city" },
      { type: "location", value: "ward" },
    ],
  });
});

test("getRawDataCatalogTree by cities", () => {
  expect(getRawDataCatalogTree(dataCatalog, "city", false, "")).toEqual([
    {
      id: "node-0",
      name: "全球データ",
      desc: "",
      children: [zenkyuData],
    },
    {
      id: "node-2",
      name: "東京都",
      desc: "",
      children: [
        {
          id: "node-3",
          name: "東京都23区",
          desc: "",
          children: [
            {
              id: "node-4",
              name: "千代田区",
              desc: "",
              children: [chiyodakuBldg, chiyodakuShelter],
            },
            {
              id: "node-7",
              name: "世田谷区",
              desc: "",
              children: [setagayakuBldg, setagayakuShelter],
            },
            tokyo23kuPark,
          ],
        },
        {
          id: "node-11",
          name: "八王子市",
          desc: "",
          children: [hachiojiBldg, hachiojiLandmark],
        },
        {
          id: "node-14",
          name: "ユースケース",
          desc: "",
          children: [tokyoUsecase],
        },
      ],
    },
    {
      id: "node-16",
      name: "栃木県",
      desc: "",
      children: [
        {
          id: "node-17",
          name: "宇都宮市",
          desc: "",
          children: [
            utsunomiyashiBldg,
            {
              id: "node-19",
              name: "都市計画決定情報モデル",
              desc: "",
              children: [utsunomiyashiUseDictrict],
            },
          ],
        },
      ],
    },
  ]);
});

test("getRawDataCatalogTree by cities as custom dataset", () => {
  expect(getRawDataCatalogTree(dataCatalog, "city", true, "")).toEqual([
    {
      id: "node-0",
      name: "全球データ",
      desc: "",
      children: [zenkyuData],
    },
    {
      id: "node-2",
      name: "東京都",
      desc: "",
      children: [
        {
          id: "node-3",
          name: "東京都23区",
          desc: "",
          children: [
            {
              id: "node-4",
              name: "千代田区",
              desc: "",
              children: [chiyodakuBldg, chiyodakuShelter],
            },
            {
              id: "node-7",
              name: "世田谷区",
              desc: "",
              children: [setagayakuBldg, setagayakuShelter],
            },
            tokyo23kuPark,
          ],
        },
        {
          id: "node-11",
          name: "八王子市",
          desc: "",
          children: [hachiojiBldg, hachiojiLandmark],
        },
        {
          id: "tokyoUsecase",
          name: "usecase",
          path: ["東京都", "usecase"],
          pref: "東京都",
          pref_code: "13",
          pref_code_i: 13,
          city_code_i: NaN,
          ward_code_i: NaN,
          type: "ユースケース",
          type_en: "usecase",
          code: 13000,
          format: "",
          url: "",
          desc: "",
          year: 2021,
          root_type: true,
        },
      ],
    },
    {
      id: "node-15",
      name: "栃木県",
      desc: "",
      children: [
        {
          id: "node-16",
          name: "宇都宮市",
          desc: "",
          children: [
            utsunomiyashiBldg,
            {
              id: "h",
              type: "都市計画決定情報モデル",
              type_en: "urf",
              type2: "用途地域",
              type2_en: "UseDistrict",
              name: "用途地域（宇都宮市）",
              path: ["栃木県", "宇都宮市", "用途地域（宇都宮市）"],
              pref: "栃木県",
              pref_code: "09",
              pref_code_i: 9,
              city: "宇都宮市",
              city_en: "utsunomiya-shi",
              city_code: "09201",
              city_code_i: 9201,
              ward_code_i: NaN,
              code: 9201,
              format: "",
              url: "",
              desc: "",
              year: 2022,
              root_type: true,
            },
          ],
        },
      ],
    },
  ]);
});

test("getRawDataCatalogTree by types", () => {
  expect(getRawDataCatalogTree(dataCatalog, "type", false, "")).toEqual([
    {
      id: "node-0",
      name: "建築物モデル",
      desc: "",
      children: [
        {
          id: "node-1",
          name: "東京都",
          desc: "",
          children: [
            {
              id: "node-2",
              name: "東京都23区",
              desc: "",
              children: [chiyodakuBldgByType, setagayakuBldgByType],
            },
            hachiojiBldgByType,
          ],
        },
        {
          id: "node-6",
          name: "栃木県",
          desc: "",
          children: [utsunomiyashiBldgByType],
        },
      ],
    },
    {
      id: "node-8",
      name: "都市計画決定情報モデル",
      desc: "",
      children: [
        {
          id: "node-9",
          name: "栃木県",
          desc: "",
          children: [
            {
              id: "node-10",
              name: "宇都宮市",
              desc: "",
              children: [utsunomiyashiUseDictrictByType],
            },
          ],
        },
      ],
    },
    {
      id: "node-12",
      name: "避難施設情報",
      desc: "",
      children: [
        {
          id: "node-13",
          name: "東京都",
          desc: "",
          children: [
            {
              id: "node-14",
              name: "東京都23区",
              desc: "",
              children: [chiyodakuShelterByType, setagayakuShelterByType],
            },
          ],
        },
      ],
    },
    {
      id: "node-17",
      name: "ランドマーク情報",
      desc: "",
      children: [
        {
          id: "node-18",
          name: "東京都",
          desc: "",
          children: [hachiojiLandmarkByType],
        },
      ],
    },
    {
      id: "node-20",
      name: "公園情報",
      desc: "",
      children: [
        {
          id: "node-21",
          name: "東京都",
          desc: "",
          children: [tokyo23kuParkByType],
        },
      ],
    },
    {
      id: "node-23",
      name: "ユースケース",
      desc: "",
      children: [
        {
          id: "node-24",
          name: "全球データ",
          desc: "",
          children: [zenkyuDataByType],
        },
        {
          id: "node-26",
          name: "東京都",
          desc: "",
          children: [tokyoUsecaseByType],
        },
      ],
    },
  ]);
});

test("getRawDataCatalogTree filter", () => {
  expect(getRawDataCatalogTree(dataCatalog, "type", false, "世田谷")).toEqual([
    {
      id: "node-0",
      name: "建築物モデル",
      desc: "",
      children: [
        {
          id: "node-1",
          name: "東京都",
          desc: "",
          children: [
            {
              id: "node-2",
              name: "東京都23区",
              desc: "",
              children: [setagayakuBldgByType],
            },
          ],
        },
      ],
    },
    {
      id: "node-4",
      name: "避難施設情報",
      desc: "",
      children: [
        {
          id: "node-5",
          name: "東京都",
          desc: "",
          children: [
            {
              id: "node-6",
              name: "東京都23区",
              desc: "",
              children: [setagayakuShelterByType],
            },
          ],
        },
      ],
    },
  ]);
});

const chiyodakuBldg: RawDataCatalogItem = {
  id: "a",
  type: "建築物モデル",
  type_en: "bldg",
  name: "建築物モデル（千代田区）",
  path: ["東京都", "東京都23区", "千代田区", "建築物モデル（千代田区）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "東京都23区",
  city_en: "tokyo-23ku",
  city_code: "13100",
  city_code_i: 13100,
  ward: "千代田区",
  ward_en: "chiyoda-ku",
  ward_code: "13101",
  ward_code_i: 13101,
  code: 13101,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const chiyodakuBldgByType = {
  ...chiyodakuBldg,
  path: ["建築物モデル", "東京都", "東京都23区", "建築物モデル（千代田区）"],
};

const chiyodakuShelter: RawDataCatalogItem = {
  id: "b",
  type: "避難施設情報",
  type_en: "shelter",
  name: "避難施設情報（千代田区）",
  path: ["東京都", "東京都23区", "千代田区", "避難施設情報（千代田区）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "東京都23区",
  city_en: "tokyo-23ku",
  city_code: "13100",
  city_code_i: 13100,
  ward: "千代田区",
  ward_en: "chiyoda-ku",
  ward_code: "13101",
  ward_code_i: 13101,
  code: 13101,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const chiyodakuShelterByType = {
  ...chiyodakuShelter,
  path: ["避難施設情報", "東京都", "東京都23区", "避難施設情報（千代田区）"],
};

const setagayakuBldg: RawDataCatalogItem = {
  id: "c",
  type: "建築物モデル",
  type_en: "bldg",
  name: "建築物モデル（世田谷区）",
  path: ["東京都", "東京都23区", "世田谷区", "建築物モデル（世田谷区）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "東京都23区",
  city_en: "tokyo-23ku",
  city_code: "13100",
  city_code_i: 13100,
  ward: "世田谷区",
  ward_en: "setagaya-ku",
  ward_code: "13112",
  ward_code_i: 13112,
  code: 13112,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const setagayakuBldgByType = {
  ...setagayakuBldg,
  path: ["建築物モデル", "東京都", "東京都23区", "建築物モデル（世田谷区）"],
};

const setagayakuShelter: RawDataCatalogItem = {
  id: "d",
  type: "避難施設情報",
  type_en: "shelter",
  name: "避難施設情報（世田谷区）",
  path: ["東京都", "東京都23区", "世田谷区", "避難施設情報（世田谷区）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "東京都23区",
  city_en: "tokyo-23ku",
  city_code: "13100",
  city_code_i: 13100,
  ward: "世田谷区",
  ward_en: "setagaya-ku",
  ward_code: "13112",
  ward_code_i: 13112,
  code: 13112,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const setagayakuShelterByType = {
  ...setagayakuShelter,
  path: ["避難施設情報", "東京都", "東京都23区", "避難施設情報（世田谷区）"],
};

const tokyo23kuPark: RawDataCatalogItem = {
  id: "e",
  type: "公園情報",
  type_en: "park",
  name: "公園情報（東京都23区）",
  path: ["東京都", "東京都23区", "公園情報（東京都23区）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "東京都23区",
  city_en: "tokyo-23ku",
  city_code: "13100",
  city_code_i: 13100,
  ward_code_i: NaN,
  code: 13100,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const tokyo23kuParkByType = {
  ...tokyo23kuPark,
  path: ["公園情報", "東京都", "公園情報（東京都23区）"],
};

const hachiojiBldg: RawDataCatalogItem = {
  id: "f",
  type: "建築物モデル",
  type_en: "bldg",
  name: "建築物モデル（八王子市）",
  path: ["東京都", "八王子市", "建築物モデル（八王子市）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "八王子市",
  city_en: "hachioji-shi",
  city_code: "13201",
  city_code_i: 13201,
  ward_code_i: NaN,
  code: 13201,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const hachiojiBldgByType = {
  ...hachiojiBldg,
  path: ["建築物モデル", "東京都", "建築物モデル（八王子市）"],
};

const hachiojiLandmark: RawDataCatalogItem = {
  id: "f",
  type: "ランドマーク情報",
  type_en: "landmark",
  name: "ランドマーク情報（八王子市）",
  path: ["東京都", "八王子市", "ランドマーク情報（八王子市）"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city: "八王子市",
  city_en: "hachioji-shi",
  city_code: "13201",
  city_code_i: 13201,
  ward_code_i: NaN,
  code: 13201,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const hachiojiLandmarkByType = {
  ...hachiojiLandmark,
  path: ["ランドマーク情報", "東京都", "ランドマーク情報（八王子市）"],
};

const tokyoUsecase: RawDataCatalogItem = {
  id: "tokyoUsecase",
  name: "usecase",
  path: ["東京都", "ユースケース", "usecase"],
  pref: "東京都",
  pref_code: "13",
  pref_code_i: 13,
  city_code_i: NaN,
  ward_code_i: NaN,
  type: "ユースケース",
  type_en: "usecase",
  code: 13000,
  format: "",
  url: "",
  desc: "",
  year: 2021,
  root_type: true,
};

const tokyoUsecaseByType = {
  ...tokyoUsecase,
  path: ["ユースケース", "東京都", "usecase"],
};

const utsunomiyashiBldg: RawDataCatalogItem = {
  id: "g",
  type: "建築物モデル",
  type_en: "bldg",
  name: "建築物モデル（宇都宮市）",
  path: ["栃木県", "宇都宮市", "建築物モデル（宇都宮市）"],
  pref: "栃木県",
  pref_code: "09",
  pref_code_i: 9,
  city: "宇都宮市",
  city_en: "utsunomiya-shi",
  city_code: "09201",
  city_code_i: 9201,
  ward_code_i: NaN,
  code: 9201,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const utsunomiyashiBldgByType = {
  ...utsunomiyashiBldg,
  path: ["建築物モデル", "栃木県", "建築物モデル（宇都宮市）"],
};

const utsunomiyashiUseDictrict: RawDataCatalogItem = {
  id: "h",
  type: "都市計画決定情報モデル",
  type_en: "urf",
  type2: "用途地域",
  type2_en: "UseDistrict",
  name: "用途地域（宇都宮市）",
  path: ["栃木県", "宇都宮市", "都市計画決定情報モデル", "用途地域（宇都宮市）"],
  pref: "栃木県",
  pref_code: "09",
  pref_code_i: 9,
  city: "宇都宮市",
  city_en: "utsunomiya-shi",
  city_code: "09201",
  city_code_i: 9201,
  ward_code_i: NaN,
  code: 9201,
  format: "",
  url: "",
  desc: "",
  year: 2022,
  root_type: true,
};

const utsunomiyashiUseDictrictByType = {
  ...utsunomiyashiUseDictrict,
  path: ["都市計画決定情報モデル", "栃木県", "宇都宮市", "用途地域（宇都宮市）"],
};

const zenkyuData: RawDataCatalogItem = {
  id: "z",
  type: "ユースケース",
  type_en: "usecase",
  name: "zenkyu",
  path: ["全球データ", "zenkyu"],
  pref: "全球データ",
  pref_code_i: 0,
  city_code_i: NaN,
  ward_code_i: NaN,
  code: 0,
  format: "",
  url: "",
  desc: "",
  year: 2022,
};

const zenkyuDataByType = {
  ...zenkyuData,
  path: ["ユースケース", "全球データ", "zenkyu"],
};

const dataCatalog: RawDataCatalogItem[] = [
  utsunomiyashiBldg,
  hachiojiBldg,
  setagayakuShelter,
  chiyodakuBldg,
  zenkyuData,
  tokyo23kuPark,
  chiyodakuShelter,
  setagayakuBldg,
  hachiojiLandmark,
  tokyoUsecase,
  utsunomiyashiUseDictrict,
];
