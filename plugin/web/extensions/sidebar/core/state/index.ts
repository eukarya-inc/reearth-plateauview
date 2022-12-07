import { atom, useAtom } from "jotai";

import { ReearthApi } from "../types";

export const defaultSettings: ReearthApi = {
  default: {
    sceneMode: "3d",
    depthTestAgainstTerrain: false,
  },
  terrain: {
    terrain: true,
    terrainType: "cesiumion",
    terrainCesiumIonAsset: "1",
  },
  tiles: [
    {
      id: "tokyo",
      tile_url: "https://cyberjapandata.gsi.go.jp/xyz/pale/{z}/{x}/{y}.png",
      tile_type: "url",
    },
  ],
};

const currentOverrides = atom<ReearthApi>(defaultSettings);
export const useCurrentOverrides = () => useAtom(currentOverrides);

const publishUrl = atom<string | undefined>(undefined);
export const usePublishUrl = () => useAtom(publishUrl);

// Is this needed? How should it be comformed? Can we just use currentOverrides?
const currentProject = atom<{}>({});
export const useCurrentProject = () => useAtom(currentProject);
