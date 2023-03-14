import { generateID } from "@web/extensions/sidebar/utils";
import { fieldGroups } from "@web/extensions/sidebar/utils/fieldGroups";
import { useMemo } from "react";

import { fieldName } from "./Fields/types";

type FieldDropdownItem = {
  [key: string]: { name: string; onClick: (property: any) => void };
};

export default ({
  onFieldAdd,
}: {
  onFieldAdd: (property: any) => ({ key }: { key: string }) => void;
}) => {
  const generalFields: FieldDropdownItem = useMemo(() => {
    return {
      idealZoom: {
        name: fieldName["idealZoom"],
        onClick: onFieldAdd({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
      },
      description: {
        name: fieldName["description"],
        onClick: onFieldAdd({}),
      },
      legend: {
        name: fieldName["legend"],
        onClick: onFieldAdd({
          style: "square",
          items: [{ title: "新しいアイテム", color: "#00bebe" }],
        }),
      },
      realtime: {
        name: fieldName["realtime"],
        onClick: onFieldAdd({ updateInterval: 30, userSettings: {} }),
      },
      timeline: {
        name: fieldName["timeline"],
        onClick: onFieldAdd({ timeFieldName: "", userSettings: { timeBasedDisplay: true } }),
      },
      currentTime: {
        name: fieldName["currentTime"],
        onClick: onFieldAdd({ date: "", time: "" }),
      },
      styleCode: {
        name: fieldName["styleCode"],
        onClick: onFieldAdd({ src: " " }),
      },
      story: {
        name: fieldName["story"],
        onClick: onFieldAdd({}),
      },
      buttonLink: {
        name: fieldName["buttonLink"],
        onClick: onFieldAdd({}),
      },
      switchGroup: {
        name: fieldName["switchGroup"],
        onClick: onFieldAdd({
          title: "Switch Group",
          groups: [
            {
              id: generateID(),
              title: "新グループ1",
              fieldGroupID: fieldGroups[0].id,
              userSettings: {},
            },
          ],
        }),
      },
      switchDataset: {
        name: fieldName["switchDataset"],
        onClick: onFieldAdd({ userSettings: {} }),
      },
      switchField: {
        name: fieldName["switchField"],
        onClick: onFieldAdd({ userSettings: {} }),
      },
      template: {
        name: fieldName["template"],
        onClick: onFieldAdd({}),
      },
      eventField: {
        name: fieldName["eventField"],
        onClick: onFieldAdd({
          eventType: "select",
          triggerEvent: "openUrl",
          urlType: "manual",
        }),
      },
      infoboxStyle: {
        name: fieldName["infoboxStyle"],
        onClick: onFieldAdd({
          displayStyle: null,
        }),
      },
    };
  }, [onFieldAdd]);

  const pointFields: FieldDropdownItem = useMemo(() => {
    return {
      pointColor: {
        name: fieldName["pointColor"],
        onClick: onFieldAdd({}),
      },
      // pointColorGradient: {
      //   name: fieldName["pointColorGradient"],
      //   onClick: ({ key }) => console.log("do something: ", key),
      // },
      pointSize: {
        name: fieldName["pointSize"],
        onClick: onFieldAdd({}),
      },
      pointIcon: {
        name: fieldName["pointIcon"],
        onClick: onFieldAdd({
          size: 1,
        }),
      },
      pointLabel: {
        name: fieldName["pointLabel"],
        onClick: onFieldAdd({}),
      },
      pointModel: {
        name: fieldName["pointModel"],
        onClick: onFieldAdd({
          scale: 1,
        }),
      },
      pointStroke: {
        name: fieldName["pointStroke"],
        onClick: onFieldAdd({}),
      },
      pointCSV: {
        name: fieldName["pointCSV"],
        onClick: onFieldAdd({}),
      },
    };
  }, [onFieldAdd]);

  const polylineFields: FieldDropdownItem = useMemo(() => {
    return {
      polylineColor: {
        name: fieldName["polylineColor"],
        onClick: onFieldAdd({}),
      },
      // polylineColorGradient: {
      //   name: fieldName["polylineColorGradient"],
      //   onClick: onFieldAdd({}),
      // },
      polylineStrokeWeight: {
        name: fieldName["polylineStrokeWeight"],
        onClick: onFieldAdd({}),
      },
    };
  }, [onFieldAdd]);

  const polygonFields: FieldDropdownItem = useMemo(() => {
    return {
      polygonColor: {
        name: fieldName["polygonColor"],
        onClick: onFieldAdd({}),
      },
      // polygonColorGradient: {
      //   name: fieldName["polygonColorGradient"],
      //   onClick: ({ key }) => console.log("do something: ", key),
      // },
      polygonStroke: {
        name: fieldName["polygonStroke"],
        onClick: onFieldAdd({}),
      },
    };
  }, [onFieldAdd]);

  const ThreeDTileFields: FieldDropdownItem = useMemo(() => {
    return {
      buildingColor: {
        name: fieldName["buildingColor"],
        onClick: onFieldAdd({ userSettings: { colorType: "none" } }),
      },
      buildingFilter: {
        name: fieldName["buildingFilter"],
        onClick: onFieldAdd({ userSettings: {} }),
      },
      buildingShadow: {
        name: fieldName["buildingShadow"],
        onClick: onFieldAdd({ userSettings: { shadow: "disabled" } }),
      },
      buildingTransparency: {
        name: fieldName["buildingTransparency"],
        onClick: onFieldAdd({ userSettings: { transparency: 100 } }),
      },
      clipping: {
        name: fieldName["clipping"],
        onClick: onFieldAdd({
          userSettings: {
            enabled: false,
            show: false,
            aboveGroundOnly: false,
            direction: "inside",
          },
        }),
      },
      floodColor: {
        name: fieldName["floodColor"],
        onClick: onFieldAdd({ userSettings: { colorType: "water" } }),
      },
      floodFilter: {
        name: fieldName["floodFilter"],
        onClick: onFieldAdd({ userSettings: {} }),
      },
    };
  }, [onFieldAdd]);

  const fieldComponentsList: { [key: string]: { name: string; fields: FieldDropdownItem } } =
    useMemo(() => {
      return {
        general: { name: "一般", fields: generalFields },
        point: { name: "ポイント", fields: pointFields },
        polyline: { name: "ポリライン", fields: polylineFields },
        polygon: { name: "ポリゴン", fields: polygonFields },
        // "3d-model": { name: "3Dモデル", fields: ThreeDModelFields },
        "3d-tile": { name: "3Dタイル", fields: ThreeDTileFields },
      };
    }, [generalFields, pointFields, polygonFields, polylineFields, ThreeDTileFields]);

  return fieldComponentsList;
};

export const cleanseOverrides: { [key: string]: any } = {
  eventField: { events: undefined },
  realtime: { data: { updateInterval: undefined } },
  timeline: { data: { time: undefined } },
  infoboxStyle: { infoboxStyle: undefined },
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
};
