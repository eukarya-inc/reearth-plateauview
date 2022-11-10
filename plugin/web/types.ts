type actionType =
  | "updateOverrides"
  | "screenshot"
  | "screenshot-save"
  | "modal-open"
  | "modal-close";
export type PostMessageProps = { action: actionType; payload?: any };
