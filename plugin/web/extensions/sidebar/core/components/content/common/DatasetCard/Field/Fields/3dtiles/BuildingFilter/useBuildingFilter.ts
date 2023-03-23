import { RefObject, useCallback, useEffect, useRef } from "react";

import { BaseFieldProps } from "../../types";

import { OptionsState } from "./constants";

export const useBuildingFilter = ({
  options,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"buildingFilter">, "dataID"> & {
  options: OptionsState;
  onUpdate: (property: any) => void;
}) => {
  const onUpdateRef = useRef(onUpdate);
  useEffect(() => {
    onUpdateRef.current = onUpdate;
  }, [onUpdate]);

  const render = useCallback(async () => {
    renderTileset(
      {
        dataID,
        options,
      },
      onUpdateRef,
    );
  }, [options, dataID]);

  useEffect(() => {
    render();
  }, [render]);
};

export type State = {
  dataID: string | undefined;
  options: OptionsState;
};

const renderTileset = (state: State, onUpdateRef: RefObject<(property: any) => void>) => {
  const updateTileset = () => {
    if (!Object.keys(state.options || {}).length) {
      return;
    }

    const defaultConditionalValue = (prop: string) =>
      `((\${${prop}} === "" || \${${prop}} === null || isNaN(Number(\${${prop}}))) ? 1 : Number(\${${prop}}))`;
    const condition = (
      max: number,
      range: [from: number, to: number] | undefined,
      conditionalValue: string,
    ) =>
      max === range?.[1]
        ? `${conditionalValue} >= ${range?.[0]}`
        : `${conditionalValue} >= ${range?.[0]} && ${conditionalValue} <= ${range?.[1]}`;
    const conditions = Object.entries(state.options || {}).reduce((res, [, v]) => {
      const conditionalValue = defaultConditionalValue(v.featurePropertyName);
      const conditionDef = condition(v.max, v.value, conditionalValue);
      return `${res ? `${res} && ` : ""}${conditionDef}`;
    }, "");
    onUpdateRef.current?.({
      show: {
        expression: {
          conditions: [...(conditions ? [[conditions, "true"]] : []), ["true", "false"]],
        },
      },
    });
  };

  updateTileset();
};
