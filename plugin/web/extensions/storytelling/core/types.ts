export type ActionType =
  | "resize"
  | "captureScene"
  | "viewScene"
  | "recaptureScene"
  | "editScene"
  | "closeSceneEditor"
  | "saveScene"
  | "getViewport"
  | "shareStory"
  | "saveStoryData"
  | "cancelPlayStory";

export type PostMessageProps = { action: ActionType; payload?: any };

export type Story = {
  id?: string;
  title?: string;
  scenes: Scene[];
};

export type Scene = {
  id: string;
  title: string;
  description: string;
  camera: Camera | undefined;
};

// Reearth types
export type Camera = {
  lat: number;
  lng: number;
  height: number;
  heading: number;
  pitch: number;
  roll: number;
  fov: number;
};

export type Viewport = {
  width: number;
  height: number;
  isMobile: boolean;
};

export type PluginExtensionInstance = {
  id: string;
  pluginId: string;
  name: string;
  extensionId: string;
  extensionType: "widget" | "block";
};

// Communications
export type PluginMessage = {
  data: EditStory | SaveStory | DeleteStory | PlayStory | CancelPlayStory;
  sender: string;
};

// sidebar -> storytelling
export type EditStory = {
  action: "editStory";
  payload: {
    id: string;
    scenes: string;
    title?: string;
  };
};

export type SaveStory = {
  action: "saveStory";
  payload: {
    id: string;
  };
};

export type DeleteStory = {
  action: "deleteStory";
  payload: {
    id: string;
  };
};

export type PlayStory = {
  action: "playStory";
  payload: {
    id: string;
    scenes: string;
    title?: string;
  };
};

// sidebar -> storytelling
// storytelling -> sidebar
export type CancelPlayStory = {
  action: "cancelPlayStory";
  payload: {
    id: string;
  };
};

// storytelling -> sidebar
export type ShareStory = {
  action: "shareStory";
  payload: {
    scenes: string;
  };
};

export type SaveStoryData = {
  action: "saveStoryData";
  payload: {
    id: string;
    scenes: string;
  };
};
