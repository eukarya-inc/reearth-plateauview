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
  | "shareStory";

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

export type PluginMessage = {
  data: any;
  sender: string;
};

export type PluginExtensionInstance = {
  id: string;
  pluginId: string;
  name: string;
  extensionId: string;
  extensionType: "widget" | "block";
};
