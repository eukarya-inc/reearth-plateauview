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
  | "saveStoryData";

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

// story title is NOT in use by story telling widget
export type EditStory = {
  type: "editStory";
  payload: {
    id: string;
    scenes: string;
    title?: string;
  };
};

// story telling will carry back this id
export type SaveStory = {
  type: "saveStory";
  payload: {
    id: string;
  };
};

// story telling will clear content if current editing this story
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
  };
};

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
