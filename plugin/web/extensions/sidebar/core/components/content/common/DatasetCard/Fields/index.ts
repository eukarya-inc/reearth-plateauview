import { ComponentType } from "react";

import IdealZoom from "./IdealZoom";
import Legend from "./Legend";

export type BasicFieldProps<T = any> = {
  id: string;
  type: T;
  title: string;
  url?: string;
};

export type Component<BP = any> = ComponentType<BasicFieldProps<BP>>;

const fields: Record<string, Component> = {
  idealZoom: IdealZoom,
  legend: Legend,
};

export default fields;
