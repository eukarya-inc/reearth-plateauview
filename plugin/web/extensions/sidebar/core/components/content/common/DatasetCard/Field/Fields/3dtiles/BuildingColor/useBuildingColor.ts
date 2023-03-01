import {
  ContextValue,
  useBuildingColorContext,
} from "@web/extensions/sidebar/core/BuildingColorContext";
import { getRGBAFromString, RGBA, rgbaToString } from "@web/extensions/sidebar/utils/color";
import debounce from "lodash/debounce";
import {
  Dispatch,
  RefObject,
  SetStateAction,
  useCallback,
  useEffect,
  useMemo,
  useRef,
} from "react";

import { BaseFieldProps } from "../../types";

import { COLOR_TYPE_CONDITIONS, makeSelectedFloodCondition } from "./conditions";
import { INDEPENDENT_COLOR_TYPE } from "./constants";

export const useBuildingColor = ({
  options,
  initialized,
  floods,
  dataID,
  onUpdate,
}: Pick<BaseFieldProps<"buildingColor">, "dataID"> & {
  initialized: boolean;
  options: Omit<BaseFieldProps<"buildingColor">["value"], "id" | "group" | "type">;
  floods: { id: string; label: string; featurePropertyName: string }[];
  onUpdate: (property: any) => void;
}) => {
  const renderRef = useRef<() => void>();
  const debouncedRender = useMemo(
    () => debounce(() => renderRef.current?.(), 100, { maxWait: 300 }),
    [],
  );
  const [, setColor] = useBuildingColorContext();

  const onUpdateRef = useRef(onUpdate);
  useEffect(() => {
    onUpdateRef.current = onUpdate;
  }, [onUpdate]);

  const render = useCallback(async () => {
    renderTileset(
      {
        dataID,
        floods,
        colorType: options.colorType,
      },
      onUpdateRef,
      setColor,
    );
  }, [options.colorType, dataID, floods, setColor]);

  useEffect(() => {
    if (!initialized) {
      return;
    }
    renderRef.current = render;
    debouncedRender();
  }, [render, debouncedRender, initialized]);
};

export type State = {
  dataID: string | undefined;
  floods: { id: string; label: string; featurePropertyName: string }[];
  colorType: string;
};

const renderTileset = (
  state: State,
  onUpdateRef: RefObject<(property: any) => void>,
  setColor: Dispatch<SetStateAction<ContextValue>>,
) => {
  const updateTileset = async () => {
    setColor(color => {
      // We can get transparency from RGBA. Because the color is defined as RGBA.;
      const transparency = color.transparency;

      const expression = {
        expression: {
          conditions: (
            COLOR_TYPE_CONDITIONS[
              (state.colorType as keyof typeof INDEPENDENT_COLOR_TYPE) || "none"
            ] ??
            makeSelectedFloodCondition(
              state.floods?.find(f => f.id === state.colorType)?.featurePropertyName,
            )
          ).map(([k, v]: [string, string]) => {
            const rgba = getRGBAFromString(v);
            if (!rgba) {
              return [k, v];
            }
            const composedRGBA = [...rgba.slice(0, -1), transparency || rgba[3]] as RGBA;
            return [k, rgbaToString(composedRGBA)];
          }),
        },
      };

      onUpdateRef.current?.({
        color: expression,
      });

      return {
        ...color,
        color: expression,
      };
    });
  };

  updateTileset();
};
