export type MouseEventData = {
  x?: number;
  y?: number;
  lat?: number;
  lng?: number;
  height?: number;
  layerId?: string;
  delta?: number;
};
type actionType = "modal-google-open" | "modal-Terrain-open" | "modal-close";

export type PostMessageProps = { action: actionType; payload?: any };
