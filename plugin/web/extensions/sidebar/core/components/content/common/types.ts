// export type Dataset = {
//   id: string;
//   filename?: string;
//   name: string;
//   description?: string;
//   hidden: boolean;
//   idealZoom: Location;
//   dataUrl: string;
//   tags?: { type: "location" | "data-type"; name: string }[];
//   fields: Field[];
// };

export type Field = {
  id: string;
  title?: string;
  value?: any;
};

export type BaseField = Field & {
  icon?: string;
  onClick?: () => void;
};

export type Location = {
  lat: number;
  lon: number;
  height: number;
  // check reearth for full params
};

export type Template = {
  id: string;
  name: string;
  fields?: Field[];
};
