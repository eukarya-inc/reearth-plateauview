export type PostMessageType =
  | "resize"
  | "minimize"
  | "captureScene"
  | "viewStory"
  | "recapture"
  | "editStory"
  | "closeStoryEditor"
  | "saveStory"
  | "getViewport"
  | "shareStoryTelling";

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

export type StoryTelling = {
  id?: string;
  title?: string;
  stories: Story[];
};

export type Story = {
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
