export type MouseEventData = {
  lat?: number;
  lng?: number;  
};
type actionType = "modal-google-open" | "modal-Terrain-open" | "modal-close";

export type PostMessageProps = { action: actionType; payload?: any };
