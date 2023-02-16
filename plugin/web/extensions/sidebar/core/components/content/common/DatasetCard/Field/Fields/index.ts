import { ComponentType } from "react";

import Clipping from "./3dtiles/Clipping";
import ButtonLink from "./general/ButtonLink";
import Description from "./general/Description";
import IdealZoom from "./general/IdealZoom";
import Legend from "./general/Legend";
import Story from "./general/Story";
import SwitchGroup from "./general/SwitchGroup";
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
  [F in FC["type"]]: { Component: ComponentType<FieldComponent & any>; hasUI: boolean };
};

const fields: Fields<FieldComponent> = {
  // general
  camera: { Component: IdealZoom, hasUI: false },
  legend: { Component: Legend, hasUI: true },
  description: { Component: Description, hasUI: true },
  switchGroup: { Component: SwitchGroup, hasUI: true },
  buttonLink: { Component: ButtonLink, hasUI: true },
  story: { Component: Story, hasUI: true },
  // point
  pointColor: { Component: PointColor, hasUI: false },
  pointColorGradient: { Component: PointColorGradient, hasUI: false },
  pointSize: { Component: PointSize, hasUI: false },
  pointIcon: { Component: PointIcon, hasUI: false },
  pointLabel: { Component: PointLabel, hasUI: false },
  pointModel: { Component: PointModel, hasUI: false },
  pointStroke: { Component: PointStroke, hasUI: false },
  // polyline
  // polygon
  // 3d-model
  clipping: { Component: Clipping, hasUI: false },
  // 3d-tile
  // realtime: Realtime,
  // template: Template,
};

export default fields;
