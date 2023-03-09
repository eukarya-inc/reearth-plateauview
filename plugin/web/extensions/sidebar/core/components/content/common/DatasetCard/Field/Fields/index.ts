import { ComponentType } from "react";

import BuildingColor from "./3dtiles/BuildingColor";
import BuildingFilter from "./3dtiles/BuildingFilter";
import BuildingShadow from "./3dtiles/BuildingShadow";
import BuildingTransparency from "./3dtiles/BuildingTransparency";
import Clipping from "./3dtiles/Clipping";
import FloodColor from "./3dtiles/FloodColor";
import FloodFilter from "./3dtiles/FloodFilter";
import ButtonLink from "./general/ButtonLink";
import CurrentTime from "./general/CurrentTime";
import Description from "./general/Description";
import EventField from "./general/EventField";
import IdealZoom from "./general/IdealZoom";
import Legend from "./general/Legend";
import Realtime from "./general/Realtime";
import Story from "./general/Story";
import StyleCode from "./general/StyleCode";
import SwitchDataset from "./general/SwitchDataset";
import SwitchGroup from "./general/SwitchGroup";
import Template from "./general/Template";
import Timeline from "./general/Timeline";
import PointColor from "./point/PointColor";
import PointColorGradient from "./point/PointColorGradient";
import PointCSV from "./point/PointCSV";
import PointIcon from "./point/PointIcon";
import PointLabel from "./point/PointLabel";
import PointModel from "./point/PointModel";
import PointSize from "./point/PointSize";
import PointStroke from "./point/PointStroke";
import PolygonColor from "./polygon/PolygonColor";
import PolygonColorGradient from "./polygon/PolygonColorGradient";
import PolygonStroke from "./polygon/PolygonStroke";
import PolylineColor from "./polyline/PolylineColor";
import PolylineColorGradient from "./polyline/PolylineColorGradient";
import PolylineStrokeWeight from "./polyline/PolylineStrokeWeight";
import { FieldComponent } from "./types";

export type Fields<FC extends FieldComponent> = {
  [F in FC["type"]]: { Component: ComponentType<FieldComponent & any>; hasUI: boolean } | null;
};

const fields: Fields<FieldComponent> = {
  // general
  idealZoom: { Component: IdealZoom, hasUI: false },
  legend: { Component: Legend, hasUI: true },
  description: { Component: Description, hasUI: true },
  switchGroup: { Component: SwitchGroup, hasUI: true },
  buttonLink: { Component: ButtonLink, hasUI: true },
  story: { Component: Story, hasUI: true },
  styleCode: { Component: StyleCode, hasUI: false },
  realtime: { Component: Realtime, hasUI: true },
  timeline: { Component: Timeline, hasUI: true },
  currentTime: { Component: CurrentTime, hasUI: false },
  switchDataset: { Component: SwitchDataset, hasUI: true },
  eventField: { Component: EventField, hasUI: false },
  // point
  pointColor: { Component: PointColor, hasUI: false },
  pointColorGradient: { Component: PointColorGradient, hasUI: false },
  pointSize: { Component: PointSize, hasUI: false },
  pointIcon: { Component: PointIcon, hasUI: false },
  pointLabel: { Component: PointLabel, hasUI: false },
  pointModel: { Component: PointModel, hasUI: false },
  pointStroke: { Component: PointStroke, hasUI: false },
  pointCSV: { Component: PointCSV, hasUI: false },
  // polyline
  polylineColor: { Component: PolylineColor, hasUI: false },
  polylineColorGradient: { Component: PolylineColorGradient, hasUI: false },
  polylineStrokeWeight: { Component: PolylineStrokeWeight, hasUI: false },
  // polygon
  polygonColor: { Component: PolygonColor, hasUI: false },
  polygonColorGradient: { Component: PolygonColorGradient, hasUI: false },
  polygonStroke: { Component: PolygonStroke, hasUI: false },
  // 3d-tile
  clipping: { Component: Clipping, hasUI: true },
  buildingFilter: { Component: BuildingFilter, hasUI: true },
  buildingTransparency: { Component: BuildingTransparency, hasUI: true },
  buildingColor: { Component: BuildingColor, hasUI: true },
  buildingShadow: { Component: BuildingShadow, hasUI: true },
  floodColor: { Component: FloodColor, hasUI: true },
  floodFilter: { Component: FloodFilter, hasUI: true },
  // 3d-model
  template: { Component: Template, hasUI: true },
};

export default fields;
