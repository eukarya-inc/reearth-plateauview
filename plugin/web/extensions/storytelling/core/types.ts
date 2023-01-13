export type PostMessageType =
  | "resize"
  | "minimize"
  | "captureScene"
  | "viewStory"
  | "recapture"
  | "editStory"
  | "closeStoryEditor"
  | "saveStory";

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
