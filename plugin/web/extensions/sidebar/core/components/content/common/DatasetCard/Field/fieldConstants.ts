export const getDefaultUserSettings = (type: any) => {
  return {
    userSettings: defaultUserSettings[type] ?? {},
  };
};

export const defaultUserSettings: { [key: string]: any } = {
  realtime: { enabled: true },
  timeline: { timeBasedDisplay: true },
  buildingColor: { colorType: "none" },
  buildingShadow: { shadow: "disabled" },
  buildingTransparency: { transparency: 100 },
  floodColor: { colorType: "water" },
};

export const cleanseOverrides: { [key: string]: any } = {
  eventField: { events: undefined },
  realtime: { data: { updateInterval: undefined } },
  timeline: { data: { time: undefined } },
  infoboxStyle: {
    infobox: {
      property: {
        default: {
          defaultContent: null,
        },
      },
    },
  },
  pointSize: { marker: { pointSize: 10 } },
  pointColor: { marker: { pointColor: "white" } },
  pointIcon: {
    marker: {
      style: "point",
      image: undefined,
      imageSize: undefined,
      imageSizeInMeters: undefined,
    },
  },
  pointLabel: {
    marker: {
      label: undefined,
      labelTypography: undefined,
      heightReference: undefined,
      labelText: undefined,
      extrude: undefined,
      labelBackground: undefined,
      labelBackgroundColor: undefined,
    },
  },
  pointModel: { model: undefined },
  pointStroke: {
    marker: {
      pointOutlineColor: undefined,
      pointOutlineWidth: undefined,
    },
  },
  polylineColor: {
    polyline: {
      strokeColor: "white",
    },
  },
  polylineStroke: {
    polyline: {
      strokeWidth: 5,
    },
  },
  polygonColor: {
    polygon: {
      fill: false,
    },
  },
  polygonStroke: {
    polygon: {
      stroke: true,
      strokeColor: "white",
      strokeWidth: 5,
    },
  },
  buildingColor: {
    "3dtiles": {
      color: "white",
    },
  },
  buildingTransparency: {
    "3dtiles": {
      color: undefined,
    },
  },
  buildingFilter: {
    "3dtiles": {
      show: true,
    },
  },
  buildingShadow: {
    "3dtiles": {
      shadows: "enabled",
    },
  },
  clipping: {
    box: undefined,
    "3dtiles": {
      experimental_clipping: undefined,
    },
  },
  floodColor: {
    "3dtiles": {
      color: undefined,
    },
  },
  floodFilter: {
    "3dtiles": {
      show: true,
    },
  },
  search: {
    "3dtiles": {
      show: true,
      color: "white",
    },
  },
};
