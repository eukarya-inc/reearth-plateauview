import { ComponentType } from "react";

import Camera from "./Camera";
import Legend from "./Legend";

export type BasicFieldProps<T = any> = {
  id: string;
  type: T;
  title: string;
  url?: string;
};

export type Component<BP = any> = ComponentType<BasicFieldProps<BP>>;

const fields: Record<string, Component> = {
  camera: Camera,
  legend: Legend,
};

export default fields;
