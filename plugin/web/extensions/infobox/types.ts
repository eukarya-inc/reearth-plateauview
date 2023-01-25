export type ActionType = "getInEditor";

export type PostMessageProps = { action: ActionType; payload?: any };

export type Primitive = {
  id: string;
  type: TilesType;
  properties: PrimitiveProperty[];
};

export type PrimitiveProperty = { key: string; value?: any };

export type TilesType = "building" | "bridge";
export type TilesTypeTitle = "建物情報" | "ブリッジ情報";

export type PublicSetting = {
  type: TilesType;
  typeTitle: TilesTypeTitle;
  properties: PublicProperty[];
};

export type PublicProperty = {
  key: string;
  title?: string;
  hidden?: boolean;
};

// Communication

// Each 3d tiles will have a tilesType (set in database)
// When adding 3d tiles into the map, sidebar need to record the new layerId - tilesType.

// When 3dtiles been selected, infobox need to know
//  1) its type
//  2) the property settings for its type
// The requests are sepreated so that the property settings for a certain type can be catched to save requests.
// Since there might be multiple select (not sure), the layerIds is an array.

// export type LayerType = {
//   layerId: string;
//   tilesType: TilesType;
// };
// // infobox -> sidebar
// export type Request3DTilesType = {
//   action: "request3DTilesType";
//   payload: {
//     layerIds: string[];
//   };
// };

// // sidebar -> infobox
// export type Get3DTilesType = {
//   action: "Get3DTilesType";
//   payload: {
//     layerTypes: LayerType[];
//   };
// };

// // infobox -> sidebar
// export type Request3DTilesPropertiesByTypes = {
//   action: "request3DTilesPropertiesByTypes";
//   payload: {
//     types: string[];
//   };
// };

// // sidebar -> infobox
// export type Get3DTilesPropertiesByTypes = {
//   action: "get3DTilesPropertiesByTypes";
//   payload: {
//     typeProperties: TypeProperties[];
//   };
// };
