// Plugin state VS Published project state (going to backend for saving)

// Or should we have state management, like Jotai, keeping track of everything separately,
// with everything coming together only on publish/share

// Main api state tree

export type API = {
  default: {
    terrain: boolean;
    sceneMode: "3d" | "2d";
  };
  tiles: Tile[];
};

type Tile = {
  id: string;
  tile_url: string;
};
