type TileSelection = "tokyo" | "gsi" | "gsi-std" | "gsi-blank";

export type ViewSelection = "3d-terrain" | "3d-smooth" | "2d";

export type BaseMapData = {
  key: TileSelection;
  tile_type?: string;
  url?: string[];
  title?: string;
  icon?: string;
};

export type MapViewData = {
  key: ViewSelection;
  title: string;
};
