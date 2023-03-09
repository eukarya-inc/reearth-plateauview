export const INDEPENDENT_COLOR_TYPE = {
  height: {
    id: "height",
    label: "高さによる塗分け",
    featurePropertyName: "計測高さ",
  },
  purpose: {
    id: "purpose",
    label: "用途による塗分け",
    featurePropertyName: "用途",
  },
  structure: {
    id: "structure",
    label: "建物構造による塗分け",
    featurePropertyName: "建物構造",
  },
  structureType: {
    id: "structureType",
    label: "構造種別による塗分け",
    featurePropertyName: "構造種別",
  },
  fireproof: {
    id: "fireproof",
    label: "耐火構造種別による塗分け",
    featurePropertyName: "建物利用現況_耐火構造種別",
  },
};

export const LEGEND_IMAGES: Record<"floods", string> = {
  floods: "https://d2jfi34fqvxlsc.cloudfront.net/main/legends/waterfloodrank/2.png",
};
