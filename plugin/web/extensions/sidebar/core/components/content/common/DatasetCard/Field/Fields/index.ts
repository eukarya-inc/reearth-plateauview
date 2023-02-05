import { ComponentType } from "react";

<<<<<<< HEAD
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
=======
// general
import Description from "./general/Description";
import IdealZoom from "./general/IdealZoom";
import Legend from "./general/Legend";
// point
import PointColor from "./point/PointColor";
import PointColorGradient from "./point/PointColorGradient";
import PointIcon from "./point/PointIcon";
import PointLabel from "./point/PointLabel";
import PointModel from "./point/PointModel";
import PointSize from "./point/PointSize";
import PointStroke from "./point/PointStroke";
import { FieldComponent } from "./types";
// import Template from "./Template";

export type Fields<FC extends FieldComponent> = {
  [F in FC["type"]]: ComponentType<FieldComponent & any>;
>>>>>>> 84df32584167e469611505ab19f7c26431abccd2
};

const fieldComponents: Fields<FieldComponent> = {
  // general
  camera: IdealZoom,
  legend: Legend,
  description: Description,
<<<<<<< HEAD
  template: Template,
  switchGroup: SwitchGroup,
=======
  // point
  pointColor: PointColor,
  pointColorGradient: PointColorGradient,
  pointSize: PointSize,
  pointIcon: PointIcon,
  pointLabel: PointLabel,
  pointModel: PointModel,
  pointStroke: PointStroke,
  // polyline
  // polygon
  // 3d-model
  // 3d-tile
  // realtime: Realtime,
  // template: Template,
>>>>>>> 84df32584167e469611505ab19f7c26431abccd2
};

export default fieldComponents;
