import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useRef } from "react";

import { BaseFieldProps } from "../../types";

export const useClippingBox = ({
  value,
  dataID,
}: Pick<BaseFieldProps<"clipping">, "value" | "dataID">) => {
  const renderer = useRef<ClippingBoxRenderer>();

  const findLayerIdFromDataset = useCallback(() => {
    return new Promise<string | undefined>(resolve => {
      const eventListenerCallback = (e: MessageEvent<any>) => {
        if (e.source !== parent) return;
        if (e.data.action === "findLayerIdFromAddedDataset") {
          resolve(e.data.payload.layerId as string);
          removeEventListener("message", eventListenerCallback);
        }
        resolve(undefined);
      };
      addEventListener("message", eventListenerCallback);
      postMsg({ action: "findLayerIdFromAddedDataset", payload: { dataID } });
    });
  }, [dataID]);

  useEffect(() => {
    const render = async () => {
      const tilesetLayerId = await findLayerIdFromDataset();
      if (value.enabled && tilesetLayerId && !renderer.current) {
        renderer.current = mountClippingBox({
          tilesetLayerId,
          keepBoxAboveGround: value.aboveGroundOnly,
          show: value.show,
          direction: value.direction,
        });
      }
      if (tilesetLayerId && renderer.current) {
        renderer.current.update({
          tilesetLayerId,
          keepBoxAboveGround: value.aboveGroundOnly,
          show: value.show,
          direction: value.direction,
        });
      }
      if ((!value.enabled || !tilesetLayerId) && renderer.current) {
        renderer.current.unmount();
        renderer.current = undefined;
      }
    };
    render();
  }, [value, findLayerIdFromDataset]);
};

const reearth = (globalThis.parent as any).reearth;

type LatLngHeight = {
  lng?: number;
  lat?: number;
  height?: number;
};

type BoxState = {
  activeBox?: boolean;
  activeScalePointIndex?: number; // 0 ~ 11
  isScalePointClicked?: boolean;
  activeEdgeIndex?: number; // 0 ~ 4
  isEdgeClicked: boolean;
  cursor?: string;
};

type ClippingBoxState = {
  tilesetLayerId: string;
  keepBoxAboveGround: boolean;
  direction: "inside" | "outside";
  show: boolean;
};

type ClippingBoxRenderer = {
  update: (state: ClippingBoxState) => void;
  unmount: () => void;
};

const mountClippingBox = (initialState: ClippingBoxState): ClippingBoxRenderer => {
  const SIDE_PLANES = [
    {
      normal: {
        x: 0,
        y: 0,
        z: 1,
      },
      distance: 0.5,
    },
    {
      normal: {
        x: 0,
        y: 0,
        z: -1,
      },
      distance: 0.5,
    },
    {
      normal: {
        x: 0,
        y: 1,
        z: 0,
      },
      distance: 0.5,
    },
    {
      normal: {
        x: 0,
        y: -1,
        z: 0,
      },
      distance: 0.5,
    },
    {
      normal: {
        x: 1,
        y: 0,
        z: 0,
      },
      distance: 0.5,
    },
    {
      normal: {
        x: -1,
        y: 0,
        z: 0,
      },
      distance: 0.5,
    },
  ];

  const viewport = reearth.camera.viewport;
  const location: LatLngHeight = {
    lng: (viewport.east + viewport.west) / 2,
    lat: (viewport.north + viewport.south) / 2,
    height: 0,
  };
  const dimensions = {
    width: 100,
    height: 100,
    length: 100,
  };
  // const clipping = {};
  const box = {
    outlineColor: "#ffffff",
    activeOutlineColor: "#0ee1ff",
    outlineWidth: 1,
    draggableOutlineWidth: 10,
    draggableOutlineColor: "rgba(14, 225, 255, 0.5)",
    activeDraggableOutlineColor: "rgba(14, 225, 255, 1)",
    fillColor: "rgba(255, 255, 255, 0.1)",
    axisLineColor: "#ffffff",
    axisLineWidth: "#ffffff",
    pointFillColor: "rgba(255, 255, 255, 0.3)",
    pointOutlineColor: "rgba(14, 225, 255, 0.5)",
    activePointOutlineColor: "rgba(14, 225, 255, 1)",
    pointOutlineWidth: 1,
  };

  const state = {
    keepBoxAboveGround: !!initialState.keepBoxAboveGround,
    tilesetLayerId: initialState.tilesetLayerId,
    isVisible: !!initialState.show,
    direction: initialState.direction,
  };

  const boxProperties: any = {
    ...dimensions,
    ...box,
  };

  const boxId = reearth.layers.add({
    type: "simple",
    data: {
      type: "geojson",
      value: {
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [location.lng, location.lat, location.height],
        },
      },
    },
    visible: state.isVisible,
    box: {
      ...boxProperties,
    },
  });
  reearth.layers.override(state.tilesetLayerId, {
    "3dtiles": {
      experimental_clipping: {
        planes: SIDE_PLANES,
        ...boxProperties,
        location: { ...location },
        visible: state.isVisible,
        direction: state.direction,
      },
    },
  });

  const lookAt = (position: LatLngHeight | null) => {
    reearth.camera.lookAt(position, { animation: false });
  };

  const allowEnterGround = () =>
    !!reearth.scene.property.default.allowEnterGround || !state.keepBoxAboveGround;

  let isBoxClicked = false;
  let isTopBottomSidePlaneClicked = false;
  let currentCameraPosition: LatLngHeight | null = null;
  let prevY: number | null = null;

  const boxState: BoxState = {
    activeBox: false,
    activeScalePointIndex: undefined, // 0 ~ 11
    isScalePointClicked: false,
    activeEdgeIndex: undefined, // 0 ~ 4
    isEdgeClicked: false,
    cursor: "default", // grab | grabbing | default
  };

  const updateBox = (shouldUpdateClipping?: boolean) => {
    reearth.layers.override(boxId, {
      visible: state.isVisible,
      data: {
        type: "geojson",
        value: {
          type: "Feature",
          geometry: {
            type: "Point",
            coordinates: [location.lng, location.lat, location.height],
          },
        },
      },
      box: {
        ...boxProperties,
        cursor: boxState.cursor,
        activeBox: boxState.activeBox,
        activeScalePointIndex: boxState.activeScalePointIndex,
        activeEdgeIndex: boxState.activeEdgeIndex,
      },
    });

    if (shouldUpdateClipping) {
      new Promise(resolve => {
        reearth.layers.override(state.tilesetLayerId, {
          "3dtiles": {
            experimental_clipping: {
              planes: SIDE_PLANES,
              ...boxProperties,
              location: { ...location },
              visible: state.isVisible,
              direction: state.direction,
            },
          },
        });
        resolve(undefined);
      });
    }
  };

  reearth.on("mousedown", (e: any) => {
    // Handle scale box
    if (e.layerId?.startsWith(`${boxId}-scale-point`)) {
      boxState.cursor = "nesw-resize";
      const index = Number(e.layerId.split("-").slice(-1)[0]);
      boxState.activeScalePointIndex = index;
      boxState.isScalePointClicked = true;
      updateBox();
    }
    // Handle edge
    if (e.layerId?.startsWith(`${boxId}-edge-draggable`)) {
      boxState.cursor = "grabbing";
      const index = Number(e.layerId.split("-").slice(-1)[0]);
      boxState.activeEdgeIndex = index;
      boxState.isEdgeClicked = true;
      updateBox();
    }

    if (e.layerId?.startsWith(`${boxId}-plane`)) {
      isBoxClicked = true;
      isTopBottomSidePlaneClicked = e.layerId.endsWith("top") || e.layerId.endsWith("bottom");
    }
    if (isBoxClicked) {
      const cameraPosition = reearth.camera.position;
      currentCameraPosition = { ...cameraPosition };
      lookAt(currentCameraPosition);

      if (!boxState.isScalePointClicked || !boxState.isEdgeClicked) {
        boxState.cursor = "grabbing";
        boxState.activeBox = true;
        updateBox();
      }
    }
  });
  reearth.on("mouseup", () => {
    if (boxState.activeScalePointIndex || boxState.activeEdgeIndex) {
      boxState.cursor = "default";

      // Handle scale box
      boxState.activeScalePointIndex = undefined;
      boxState.isScalePointClicked = false;
      // Handle edge
      boxState.activeEdgeIndex = undefined;
      boxState.isEdgeClicked = false;

      updateBox();
    }

    if (isBoxClicked) {
      // TODO: Fix to use `animation: false`.
      // This is workaround because if we use `lookAt` with `animation: false`, zooming interaction is freeze.
      reearth.camera.lookAt(currentCameraPosition, { duration: 0 });
      currentCameraPosition = null;
      isBoxClicked = false;
      isTopBottomSidePlaneClicked = false;
      prevY = null;

      boxState.activeBox = false;
      boxState.cursor = "default";
      updateBox();
    }
  });
  reearth.on("mousemove", async (e: any) => {
    if (!isBoxClicked) return;
    if (!prevY) {
      prevY = e.y;
    }

    if (isTopBottomSidePlaneClicked) {
      location.height = await (async () => {
        if (!allowEnterGround()) {
          const boxBottomHeight = (location.height || 0) - boxProperties.height / 2;
          const floorHeight = await reearth.scene.sampleTerrainHeight(location.lat, location.lng);
          if (boxBottomHeight < floorHeight) {
            return boxBottomHeight + (floorHeight - boxBottomHeight);
          }
          return location.height;
        }

        const scale =
          Math.floor(location.height || 0) > 5
            ? reearth.camera.position.height / (location.height || 0)
            : 1;
        return Math.max((location.height || 0) + ((prevY || 0) - e.y) * scale, 1);
      })();
      prevY = e.y;
    } else {
      location.lat = e.lat;
      location.lng = e.lng;
    }

    lookAt(currentCameraPosition);

    updateBox(true);
  });
  reearth.on("mouseenter", (e: any) => {
    const enableEnterHandling =
      !boxState.isScalePointClicked && !boxState.isEdgeClicked && !isBoxClicked;
    // Handle scale box
    if (e.layerId?.startsWith(`${boxId}-scale-point`)) {
      if (enableEnterHandling) {
        boxState.cursor = "nesw-resize";
        const index = Number(e.layerId.split("-").slice(-1)[0]);
        boxState.activeScalePointIndex = index;
        updateBox();
      }
    }
    // Handle edge
    if (e.layerId?.startsWith(`${boxId}-edge-draggable`)) {
      if (enableEnterHandling) {
        boxState.cursor = "grab";
        const index = Number(e.layerId.split("-").slice(-1)[0]);
        boxState.activeEdgeIndex = index;
        updateBox();
      }
    }

    if (e.layerId?.startsWith(`${boxId}-plane`)) {
      if (enableEnterHandling) {
        boxState.cursor = "grab";
        boxState.activeBox = true;
        updateBox();
      }
    }
  });
  reearth.on("mouseleave", (e: any) => {
    const enableLeaveHandling =
      !boxState.isScalePointClicked && !boxState.isEdgeClicked && !isBoxClicked;
    // Handle scale box
    if (e.layerId?.startsWith(`${boxId}-scale-point`)) {
      if (enableLeaveHandling) {
        boxState.cursor = "default";
        boxState.activeScalePointIndex = undefined;
        updateBox();
      }
    }
    // Handle edge
    if (e.layerId?.startsWith(`${boxId}-edge-draggable`)) {
      if (enableLeaveHandling) {
        boxState.cursor = "default";
        boxState.activeEdgeIndex = undefined;
        updateBox();
      }
    }

    if (e.layerId?.startsWith(`${boxId}-plane`)) {
      if (enableLeaveHandling) {
        boxState.cursor = "default";
        boxState.activeBox = false;
        updateBox();
      }
    }
  });

  reearth.on("layeredit", (e: any) => {
    if (e.layerId?.startsWith(`${boxId}-scale-point`) && e.scale) {
      lookAt(currentCameraPosition);

      const scale = e.scale;

      boxProperties.width = scale.width;
      boxProperties.height = scale.height;
      boxProperties.length = scale.length;
      location.lng = scale.location.lng;
      location.lat = scale.location.lat;
      location.height = scale.location.height;

      updateBox(true);
    }

    if (e.layerId?.startsWith(`${boxId}-edge-draggable`) && e.rotate) {
      lookAt(currentCameraPosition);

      const rotate = e.rotate;

      boxProperties.heading = rotate.heading;
      boxProperties.pitch = rotate.pitch;
      boxProperties.roll = rotate.roll;

      updateBox(true);
    }
  });

  const update = (next: ClippingBoxState) => {
    state.tilesetLayerId = next.tilesetLayerId;
    state.isVisible = next.show;
    state.keepBoxAboveGround = next.keepBoxAboveGround;
    state.direction = next.direction;
    updateBox(true);
  };
  const unmount = () => {
    reearth.layers.delete(boxId);
    reearth.layers.override(state.tilesetLayerId, {
      "3dtiles": {
        experimental_clipping: undefined,
      },
    });
  };
  return { update, unmount };
};
