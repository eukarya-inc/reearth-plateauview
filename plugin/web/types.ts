type actionType =
  | "updateOverrides"
  | "screenshot"
  | "screenshot-save"
  | "modal-open"
  | "modal-close"
  | "msgFromModal"
  | "minimize";
export type PostMessageProps = { action: actionType; payload?: any };
