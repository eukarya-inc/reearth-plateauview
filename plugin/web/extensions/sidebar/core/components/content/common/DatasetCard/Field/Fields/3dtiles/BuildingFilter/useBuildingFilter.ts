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

  const render = useCallback(async () => {
    if (!renderer.current) {
      renderer.current = mountTileset({
        dataID,
        height: value.height,
        abovegroundFloor: value.abovegroundFloor,
        basementFloor: value.basementFloor,
      });
    }
    if (renderer.current) {
      renderer.current.update({
        dataID,
        height: value.height,
        abovegroundFloor: value.abovegroundFloor,
        basementFloor: value.basementFloor,
      });
    }
  }, [value, dataID]);

  useEffect(() => {
    renderRef.current = render;
    debouncedRender();
  }, [render, debouncedRender]);
};

export type State = {
  dataID: string | undefined;
  height: [from: number, to: number];
  abovegroundFloor: [from: number, to: number];
  basementFloor: [from: number, to: number];
};

type Renderer = {
  update: (state: State) => void;
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
    postMsg({
      action: "update3dtilesShow",
      payload: {
        dataID: state.dataID,
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
  return { update };
};
