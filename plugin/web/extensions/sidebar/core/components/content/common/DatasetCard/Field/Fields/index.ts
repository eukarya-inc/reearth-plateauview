import { ComponentType } from "react";

import Description from "./Description";
import IdealZoom from "./IdealZoom";
import Legend from "./Legend";
import SwitchGroup from "./SwitchGroup";
import Template from "./Template";

export type FieldGroup = "general" | "point" | "polyline" | "polygon" | "3d-model" | "3d-tile";

export type FieldType = "template" | "idealZoom" | "description" | "legend" | "switchGroup";

export type BaseField<T extends FieldType> = {
  id: string;
  type: T;
  title: string;
  configTitle: string;
  // icon?: string;
  // url?: string;
  // onChange?: () => void;
};

export type Fields<FT extends FieldType> = {
  [F in FT]: ComponentType<BaseField<F> & any>;
};

const fields: Fields<FieldType> = {
  idealZoom: IdealZoom,
  legend: Legend,
  description: Description,
  template: Template,
  switchGroup: SwitchGroup,
};

export default fields;
