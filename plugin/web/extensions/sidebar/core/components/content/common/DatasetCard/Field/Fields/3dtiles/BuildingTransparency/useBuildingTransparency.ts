import { postMsg } from "@web/extensions/sidebar/utils";
import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useRef } from "react";

import { BaseFieldProps } from "../../types";

import { getRGBAFromString } from "./utils";

export const useBuildingTransparency = ({
  value,
  dataID,
}: Pick<BaseFieldProps<"buildingTransparency">, "value" | "dataID">) => {
  const renderer = useRef<Renderer>();
  const renderRef = useRef<() => void>();
  const debouncedRender = useMemo(
    () => debounce(() => renderRef.current?.(), 100, { maxWait: 300 }),
    [],
  );

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

  const render = useCallback(async () => {
    const layerId = await findLayerIdFromDataset();
    if (layerId && !renderer.current) {
      renderer.current = mountTileset({
        layerId,
        transparency: value.transparency,
      });
    }
    if (layerId && renderer.current) {
      renderer.current.update({
        layerId,
        transparency: value.transparency,
      });
    }
    if (!layerId && renderer.current) {
      renderer.current.unmount();
      renderer.current = undefined;
    }
  }, [value, findLayerIdFromDataset]);

  useEffect(() => {
    renderRef.current = render;
    debouncedRender();
  }, [render, debouncedRender]);
};

const reearth = (globalThis.parent as any).reearth;

export type State = {
  layerId: string;
  transparency: number;
};

type Renderer = {
  update: (state: State) => void;
  unmount: () => void;
};

const DEFAULT_RGBA = [255, 255, 255, 1];
const mountTileset = (initialState: State): Renderer => {
  const state: Partial<State> = {};
  const updateState = (next: State) => {
    Object.entries(next).forEach(([k, v]) => {
      state[k as keyof State] = v as any;
    });
  };

  let prevLayerId: string;
  let curLayer: any;
  const updateTileset = () => {
    if (prevLayerId === state.layerId) {
      curLayer = reearth.layers.findById(state.layerId);
      prevLayerId = state.layerId;
    }
    const prevRGBA = getRGBAFromString(curLayer?.["3dtiles"]?.color) || DEFAULT_RGBA;
    const prevRGBStr = prevRGBA?.slice(0, -1).join(",");
    reearth.layers.override(state.layerId, {
      "3dtiles": {
        color: `rgba(${prevRGBStr}, ${(state.transparency || 100) / 100})`,
      },
    });
  };

  // Initialize
  updateState(initialState);
  updateTileset();

  const update = (next: State) => {
    updateState(next);
    updateTileset();
  };
  const unmount = () => {
    reearth.layers.override(state.layerId, {
      "3dtiles": {
        color: undefined,
      },
    });
  };
  return { update, unmount };
};
