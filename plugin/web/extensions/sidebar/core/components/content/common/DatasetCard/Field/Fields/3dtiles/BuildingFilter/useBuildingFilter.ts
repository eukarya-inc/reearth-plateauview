import { postMsg } from "@web/extensions/sidebar/utils";
import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useRef } from "react";

import { BaseFieldProps } from "../../types";

export const MAX_HEIGHT = 200;
export const MAX_ABOVEGROUND_FLOOR = 50;
export const MAX_BASEMENT_FLOOR = 5;

export const useBuildingFilter = ({
  value,
  dataID,
}: Pick<BaseFieldProps<"buildingFilter">, "value" | "dataID">) => {
  const renderer = useRef<Renderer>();
  const renderRef = useRef<() => void>();
  const debouncedRender = useMemo(
    () => debounce(() => renderRef.current?.(), 300, { maxWait: 1000 }),
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
        height: value.height,
        abovegroundFloor: value.abovegroundFloor,
        basementFloor: value.basementFloor,
      });
    }
    if (layerId && renderer.current) {
      renderer.current.update({
        layerId,
        height: value.height,
        abovegroundFloor: value.abovegroundFloor,
        basementFloor: value.basementFloor,
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
  height: [from: number, to: number];
  abovegroundFloor: [from: number, to: number];
  basementFloor: [from: number, to: number];
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
    const defaultConditionalValue = (prop: string) =>
      `((\${${prop}} === "" || \${${prop}} === null || isNaN(Number(\${${prop}}))) ? 1 : Number(\${${prop}}))`;
    const conditionalHeight = defaultConditionalValue("計測高さ");
    const conditionalAbovegroundFloor = defaultConditionalValue("地上階数");
    const conditionalBasementFloor = defaultConditionalValue("地下階数");
    const condition = (
      max: number,
      range: [from: number, to: number] | undefined,
      conditionalValue: string,
    ) =>
      max === range?.[1]
        ? `${conditionalValue} >= ${range?.[0]}`
        : `${conditionalValue} >= ${range?.[0]} && ${conditionalValue} <= ${range?.[1]}`;

    const conditionForHeight = condition(MAX_HEIGHT, state.height, conditionalHeight);
    const conditionForAbovegroundFloor = condition(
      MAX_ABOVEGROUND_FLOOR,
      state.abovegroundFloor,
      conditionalAbovegroundFloor,
    );
    const conditionForBasementFloor = condition(
      MAX_BASEMENT_FLOOR,
      state.basementFloor,
      conditionalBasementFloor,
    );
    reearth.layers.override(state.layerId, {
      "3dtiles": {
        show: {
          expression: {
            conditions: [
              [
                `${conditionForHeight} && ${conditionForAbovegroundFloor} && ${conditionForBasementFloor}`,
                "true",
              ],
              ["true", "false"],
            ],
          },
        },
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
        show: true,
      },
    });
  };
  return { update, unmount };
};
