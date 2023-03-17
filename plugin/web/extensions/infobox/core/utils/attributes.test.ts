import { expect, test, vi } from "vitest";

import { attributesMap, getAttributes, getRootFields } from "./attributes";
import type { Json } from "./json";

test("getAttributes", () => {
  const src: Json = {
    bbb: {},
    aaa: {
      bbb: "ccc",
      ddd: [{ c: "b" }, { b: "a", a: "" }],
    },
  };
  expect(flatKeys(src)).toEqual([
    "",
    "bbb",
    "aaa",
    "aaa.bbb",
    "aaa.ddd",
    "aaa.ddd.0",
    "aaa.ddd.0.c",
    "aaa.ddd.1",
    "aaa.ddd.1.b",
    "aaa.ddd.1.a",
  ]);

  const actual = getAttributes(src);
  expect(flatKeys(actual)).toEqual([
    "",
    "AAA（aaa）",
    "AAA（aaa）.bbb",
    "AAA（aaa）.DDD（ddd）",
    "AAA（aaa）.DDD（ddd）.0",
    "AAA（aaa）.DDD（ddd）.0.c",
    "AAA（aaa）.DDD（ddd）.1",
    "AAA（aaa）.DDD（ddd）.1.a",
    "AAA（aaa）.DDD（ddd）.1.b",
    "bbb",
  ]);
});

test("getRootFields", () => {
  const res = getRootFields({
    attributes: {
      "gml:id": "id",
      "bldg:class": "堅ろう建物",
      "bldg:yearOfConstruction": 2008,
      "bldg:usage": ["共同住宅"],
      "bldg:measuredHeight": 20,
      "bldg:measuredHeightUmo": "m",
      "bldg:storeysAboveGround": 12,
      "bldg:storeysBelowGround": 0,
      "uro:majorUsage": "建物利用現況（大分類）",
      "uro:orgUsage": "建物利用現況（中分類）",
      "uro:orgUsage2": "建物利用現況（小分類）",
      "uro:detailedUsage": "建物利用現況（詳細分類）",
      "uro:BuildingDetailAttribute": [
        {
          "uro:siteArea": 100,
          "uro:buildingStructureType": "鉄筋コンクリート造",
          "uro:buildingStructureOrgType": "鉄筋コンクリート造",
          "uro:fireproofStructureType": "不明",
          "uro:surveyYear": 2018,
          "uro:urbanPlanType": "都市計画区域",
          "uro:areaClassificationType": "区域区分",
          "uro:districtsAndZonesType": ["地域地区"],
        },
      ],
      "uro:BuildingRiverFloodingRiskAttribute": [
        {
          "uro:description": "六角川水系武雄川",
          "uro:depth": 0.618,
          "uro:depth_uom": "m",
          "uro:adminType": "国",
          "uro:scale": "L2（想定最大規模）",
          "uro:rank_code": "1",
          "uro:duration": "継続時間",
        },
        {
          "uro:description": "六角川水系武雄川",
          "uro:depth": 0,
          "uro:depth_uom": "m",
          "uro:adminType": "都道府県",
          "uro:scale": "L1（計画規模）",
          "uro:rank_code": "1",
          "uro:duration": "継続時間",
        },
      ],
      "uro:BuildingLandSlideRiskAttribute": "土砂災害警戒区域",
    },
  });

  expect(res).toEqual({
    gml_id: "id",
    分類: "堅ろう建物",
    用途: "共同住宅",
    建築年: 2008,
    計測高さ: 20,
    地上階数: 12,
    地下階数: 0,
    敷地面積: 100,
    構造種別: "鉄筋コンクリート造",
    "構造種別（自治体独自）": "鉄筋コンクリート造",
    耐火構造種別: "不明",
    都市計画区域: "都市計画区域",
    区域区分: "区域区分",
    地域地区: "地域地区",
    調査年: 2018,
    "建物利用現況（大分類）": "建物利用現況（大分類）",
    "建物利用現況（中分類）": "建物利用現況（中分類）",
    "建物利用現況（小分類）": "建物利用現況（小分類）",
    "建物利用現況（詳細分類）": "建物利用現況（詳細分類）",
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_浸水ランク": "1",
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_浸水深": 0.618,
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_継続時間": "継続時間",
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_浸水ランク": "1",
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_浸水深": 0,
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_継続時間": "継続時間",
    土砂災害警戒区域: "土砂災害警戒区域",
  });

  expect(flatKeys(res)).toEqual([
    "",
    "gml_id",
    "分類",
    "用途",
    "建築年",
    "計測高さ",
    "地上階数",
    "地下階数",
    "敷地面積",
    "構造種別",
    "構造種別（自治体独自）",
    "耐火構造種別",
    "都市計画区域",
    "区域区分",
    "地域地区",
    "調査年",
    "建物利用現況（大分類）",
    "建物利用現況（中分類）",
    "建物利用現況（小分類）",
    "建物利用現況（詳細分類）",
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_浸水ランク",
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_浸水深",
    "六角川水系武雄川（国管理区間）_L2（想定最大規模）_継続時間",
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_浸水ランク",
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_浸水深",
    "六角川水系武雄川（都道府県管理区間）_L1（計画規模）_継続時間",
    "土砂災害警戒区域",
  ]);
});

test("attributesMap", () => {
  expect(attributesMap.get("ddd")).toBe("DDD");
});

function flatKeys(obj: Json, parentKey?: string): string[] {
  if (typeof obj !== "object" || !obj) return [parentKey || ""];
  return [
    parentKey || "",
    ...Object.entries(obj).flatMap(([k, v]) =>
      flatKeys(v, `${parentKey ? `${parentKey}.` : ""}${k}`),
    ),
  ];
}

vi.mock("./attributes.csv?raw", () => ({
  default: "ddd,DDD\naaa,AAA\n",
}));
