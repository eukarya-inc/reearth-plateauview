import { get } from "lodash";

import type { Properties } from "../../types";

import attributesData from "./attributes.csv?raw";
import type { Json, JsonArray, JsonObject } from "./json";

export const attributesMap = new Map<string, string>();

attributesData
  .split("\n")
  .map(l => l.split(","))
  .forEach(l => {
    if (!l || !l[0] || !l[1] || typeof l[0] !== "string" || typeof l[1] !== "string") return;
    attributesMap.set(l[0], l[1]);
  });

export function getAttributes(attributes: Json): Json {
  if (!attributes || typeof attributes !== "object") return attributes;
  return walk(attributes, attributesMap);
}

function walk(obj: JsonObject | JsonArray, keyMap?: Map<string, string>): JsonObject | JsonArray {
  if (!obj || typeof obj !== "object") return obj;

  if (Array.isArray(obj)) {
    return obj.map(o => (typeof o === "object" && o ? walk(o) : o));
  }

  return Object.fromEntries(
    Object.entries(obj)
      .sort((a, b) => a[0].localeCompare(b[0]))
      .map(([k, v]) => {
        const nk = keyMap?.get(k);
        const ak = nk ? `${nk}（${k}）` : k;

        if (typeof v === "object" && v) {
          return [ak || k, walk(v, keyMap)];
        }
        return [ak || k, v];
      }),
  );
}

export function getRootFields(properties: Properties): any {
  return filterObjects({
    gml_id: get(properties, ["attributes", "gml:id"]),
    名称: get(properties, ["attributes", "gml:name"]),
    分類: get(properties, ["attributes", "bldg:class"]),
    用途: get(properties, ["attributes", "bldg:usage", 0]),
    建築年: get(properties, ["attributes", "bldg:yearOfConstruction"]),
    計測高さ: get(properties, ["attributes", "bldg:measuredHeight"]),
    地上階数: get(properties, ["attributes", "bldg:storeysAboveGround"]),
    地下階数: get(properties, ["attributes", "bldg:storeysBelowGround"]),
    敷地面積: get(properties, ["attributes", "uro:BuildingDetailAttribute", 0, "uro:siteArea"]),
    構造種別: get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:buildingStructureType",
    ]),
    "構造種別（自治体独自）": get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:buildingStructureOrgType",
    ]),
    耐火構造種別: get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:fireproofStructureType",
    ]),
    都市計画区域: get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:urbanPlanType",
    ]),
    区域区分: get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:areaClassificationType",
    ]),
    地域地区: get(properties, [
      "attributes",
      "uro:BuildingDetailAttribute",
      0,
      "uro:districtsAndZonesType",
      0,
    ]),
    調査年: get(properties, ["attributes", "uro:BuildingDetailAttribute", 0, "uro:surveyYear"]),
    "建物利用現況（大分類）": get(properties, ["attributes", "uro:majorUsage"]),
    "建物利用現況（中分類）": get(properties, ["attributes", "uro:orgUsage"]),
    "建物利用現況（小分類）": get(properties, ["attributes", "uro:orgUsage2"]),
    "建物利用現況（詳細分類）": get(properties, ["attributes", "uro:detailedUsage"]),
    ...floodFields(properties),
    土砂災害警戒区域: get(properties, ["attributes", "uro:BuildingLandSlideRiskAttribute"]),
  });
}

function floodFields(properties: Properties): any {
  const fld = get(properties, ["attributes", "uro:BuildingRiverFloodingRiskAttribute"]) as
    | BuildingRiverFloodingRiskAttribute[]
    | undefined;
  if (!Array.isArray(fld)) return {};

  return Object.fromEntries(
    fld
      .slice(0)
      .sort(
        (a, b) =>
          a?.["uro:description"]?.localeCompare(b?.["uro:description"] || "") ||
          a?.["uro:adminType"]?.localeCompare(b?.["uro:adminType"] || "") ||
          a?.["uro:scale"]?.localeCompare(b?.["uro:scale"] || "") ||
          0,
      )
      .flatMap(a => {
        if (!a || !a["uro:description"] || !a["uro:adminType"] || !a["uro:scale"]) return [];
        const prefix = `${a["uro:description"]}（${a["uro:adminType"]}管理区間）_${a["uro:scale"]}`;
        return [
          [`${prefix}_浸水ランク`, a["uro:rank_code"]],
          [`${prefix}_浸水深`, a["uro:depth"]],
          [`${prefix}_継続時間`, a["uro:duration"]],
        ];
      })
      .filter(f => typeof f[1] !== "undefined" && (typeof f[1] !== "string" || f[1])),
  );
}

function filterObjects(obj: any): any {
  return Object.fromEntries(
    Object.entries(obj).filter(
      e => typeof e[1] !== "undefined" && (typeof e[1] !== "string" || e[1]),
    ),
  );
}

type BuildingRiverFloodingRiskAttribute = {
  "uro:description"?: string; // 指定河川名称
  "uro:depth"?: number; // 浸水深
  "uro:depth_uom"?: string; // 浸水深の単位
  "uro:rank_code"?: string; // 浸水ランクコード
  "uro:duration"?: string; // 継続時間
  "uro:adminType"?: string; //
  "uro:scale"?: string; // 浸水規模
};
