import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useRef } from "react";

import { BaseFieldProps } from "../../types";

export const MAX_HEIGHT = 200;
export const MAX_ABOVEGROUND_FLOOR = 50;
export const MAX_BASEMENT_FLOOR = 5;

export const useBuildingShadow = ({
  value,
  dataID,
}: Pick<BaseFieldProps<"buildingShadow">, "value" | "dataID">) => {
  const renderer = useRef<Renderer>();

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
        shadow: value.shadow,
      });
    }
    if (layerId && renderer.current) {
      renderer.current.update({
        layerId,
        shadow: value.shadow,
      });
    }
    if (!layerId && renderer.current) {
      renderer.current.unmount();
      renderer.current = undefined;
    }
  }, [value, findLayerIdFromDataset]);

  useEffect(() => {
    render();
  }, [render]);
};

const reearth = (globalThis.parent as any).reearth;

export type State = {
  layerId: string;
  shadow: BaseFieldProps<"buildingShadow">["value"]["shadow"];
};

type Renderer = {
  update: (state: State) => void;
  unmount: () => void;
};

const mountTileset = (initialState: State): Renderer => {
  const state: Partial<State> = {};
  const updateState = (next: State) => {
    Object.entries(next).forEach(([k, v]) => {
      state[k as keyof State] = v as any;
    });
  };

  const updateTileset = () => {
    reearth.layers.override(state.layerId, {
      "3dtiles": {
        shadows: state.shadow,
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
        shadows: undefined,
      },
    });
  };
  return { update, unmount };
};
