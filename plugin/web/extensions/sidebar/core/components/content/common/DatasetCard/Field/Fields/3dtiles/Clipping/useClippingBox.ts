import { RefObject, useEffect, useRef } from "react";

import { BaseFieldProps } from "../../types";

export const useClippingBox = ({
  options,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"clipping">, "dataID"> & {
  options: Omit<BaseFieldProps<"clipping">["value"], "id" | "group" | "type">;
  onUpdate: (tilesetProperty: any, boxProperty: any) => void;
}) => {
  const onUpdateRef = useRef(onUpdate);
  useEffect(() => {
    onUpdateRef.current = onUpdate;
  }, [onUpdate]);

  useEffect(() => {
    const render = async () => {
      await renderTileset(
        {
          dataID,
          keepBoxAboveGround: options.aboveGroundOnly,
          show: options.show,
          direction: options.direction,
          enabled: options.enabled,
        },
        onUpdateRef,
      );
    };
    render();
  }, [options.aboveGroundOnly, options.direction, options.enabled, options.show, dataID]);
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
  dataID: string | undefined;
  keepBoxAboveGround: boolean;
  direction: "inside" | "outside";
  show: boolean;
  enabled: boolean;
};

const renderTileset = async (
  state: ClippingBoxState,
  onUpdateRef: RefObject<(tilesetProperty: any, boxProperty: any) => void>,
) => {
  const viewport = reearth.viewport;
  const centerOnScreen = reearth.scene.getLocationFromScreen(
    viewport.width / 2,
    viewport.height / 2,
  );
  const dimensions = {
    width: 100,
    height: 100,
    length: 100,
  };
  const location: LatLngHeight = {
    lng: centerOnScreen.lng,
    lat: centerOnScreen.lat,
    height: dimensions.height,
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

  const boxProperties: any = {
    ...dimensions,
    ...box,
  };

  const boxState: BoxState = {
    activeBox: false,
    activeScalePointIndex: undefined, // 0 ~ 11
    isScalePointClicked: false,
    activeEdgeIndex: undefined, // 0 ~ 4
    isEdgeClicked: false,
    cursor: "default", // grab | grabbing | default
  };

  const updateBox = () => {
    onUpdateRef.current?.(
      {
        experimental_clipping: state.enabled
          ? {
              ...boxProperties,
              coordinates: [location.lng, location.lat, location.height],
              visible: state.show,
              direction: state.direction,
              allowEnterGround: !state.keepBoxAboveGround,
              useBuiltinBox: true,
            }
          : undefined,
      },
      state.enabled
        ? {
            ...boxProperties,
            cursor: boxState.cursor,
            activeBox: boxState.activeBox,
            activeScalePointIndex: boxState.activeScalePointIndex,
            activeEdgeIndex: boxState.activeEdgeIndex,
            allowEnterGround: !state.keepBoxAboveGround,
          }
        : undefined,
    );
  };

  updateBox();
};
