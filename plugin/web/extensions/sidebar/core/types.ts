export type actionType =
  | "updateOverrides"
  | "screenshot"
  | "screenshot-save"
  | "modal-open"
  | "modal-close"
  | "msgFromModal"
  | "minimize";

export type ReearthApi = {
  default?: {
    camera?: Camera;
    terrain?: boolean;
    sceneMode?: SceneMode;
    depthTestAgainstTerrain?: boolean;
    allowEnterGround?: boolean;
  };
  tiles?: Tile[];
};

export type SceneMode = "3d" | "2d";

type Tile = {
  id: string;
  tile_url: string;
  tile_type: string;
};

type Camera = {
  lat: number;
  lng: number;
  altitude: number;
  heading: number;
  pitch: number;
  roll: number;
};

export type PublishProject = {
  // Here would be all fields being saved to backend
};
