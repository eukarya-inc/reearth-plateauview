export type PostMessageType =
  | "resize"
  | "minimize"
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

export type PostMessageProps = { type: PostMessageType; payload?: any };

export type ReearthApi = {};

export type Camera = {
  lat: number;
  lng: number;
  height: number;
  heading: number;
  pitch: number;
  roll: number;
  fov: number;
};

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

export type PluginMessage = {
  data: EditStory | SaveStory | DeleteStory | PlayStory | CancelPlayStory;
  sender: string;
};

// Communications

// sidebar -> storytelling
export type EditStory = {
  type: "editStory";
  payload: {
    id: string;
    scenes: string;
    title?: string;
  };
};

export type SaveStory = {
  type: "saveStory";
  payload: {
    id: string;
  };
};

export type DeleteStory = {
  type: "deleteStory";
  payload: {
    id: string;
  };
};

export type PlayStory = {
  type: "playStory";
  payload: {
    id: string;
    scenes: string;
    title?: string;
  };
};

// sidebar -> storytelling
// storytelling -> sidebar
export type CancelPlayStory = {
  type: "cancelPlayStory";
  payload: {
    id: string;
  };
};

// storytelling -> sidebar
export type ShareStory = {
  type: "shareStory";
  payload: {
    scenes: string;
  };
};

export type SaveStoryData = {
  type: "saveStoryData";
  payload: {
    id: string;
    scenes: string;
  };
};
